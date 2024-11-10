package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	ctx := context.Background()
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	conn, err := initDB(ctx)
	if err != nil {
		e.Logger.Fatal(err)
		return
	}
	defer conn.Close(ctx)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/query-builder", func(c echo.Context) error {
		return getCars(c, conn)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func initDB(ctx context.Context) (*pgx.Conn, error) {
	host := "localhost"
	user := "sqlbuilderuser"
	password := "password"
	dbName := "sqlbuildertest"
	port := 5432

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbName, port)

	return pgx.Connect(ctx, dsn)

}

type Car struct {
	ID       int    `json:"id"`
	Brand    string `json:"brand"`
	Model    string `json:"model"`
	Year     string `json:"year"`
	State    string `json:"state"`
	Color    string `json:"color"`
	FuelType string `json:"fuel_type"`
	BodyType string `json:"body_type"`
}

func getCars(c echo.Context, conn *pgx.Conn) error {
	fullColumns := []string{"id", "brand", "model", "year", "state", "color", "fuel_type", "body_type"}
	// fields param to filter which column shown
	selectedColumns := validateFields(fullColumns, c.QueryParam("fields"))
	filterCarColumns := fullColumns[1:]

	carsSQL := sq.Select(selectedColumns...).
		PlaceholderFormat(sq.Dollar).
		From("cars")

	for _, col := range filterCarColumns {
		queryParam := c.QueryParam(col)
		if queryParam != "" {
			carsSQL = carsSQL.Where(sq.Eq{col: queryParam})
		}
	}

	// limit param to limit how much the data will be shown
	limit := validateLimit(c.QueryParam("limit"))
	carsSQL = carsSQL.Limit(limit)

	// page param to get to the desired page
	offset := validateOffset(c.QueryParam("page"), limit)
	carsSQL = carsSQL.Offset(offset)

	sortBy := "id"
	carsSQL = carsSQL.OrderBy(sortBy)

	sql, args, err := carsSQL.ToSql()
	c.Logger().Debug(sql)
	c.Logger().Debug(args)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate SQL"})
	}

	rows, err := conn.Query(c.Request().Context(), sql, args...)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch cars"})
	}
	defer rows.Close()

	var cars []map[string]interface{}

	for rows.Next() {
		car, err := scanRowToMap(rows, selectedColumns)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to scan car"})
		}

		cars = append(cars, car)
	}

	return c.JSON(http.StatusOK, cars)
}

func scanRowToMap(rows pgx.Rows, cols []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	count := len(cols)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for i := range cols {
		valuePtrs[i] = &values[i]
	}

	err := rows.Scan(valuePtrs...)
	if err != nil {
		return result, err
	}

	for i, col := range cols {
		val := values[i]

		b, ok := val.([]byte)
		if ok {
			result[col] = string(b)
		} else {
			result[col] = val
		}
	}

	return result, nil
}

func validateFields(fullFields []string, desiredFields string) []string {
	if desiredFields == "" {
		return fullFields
	}

	// Create a map for fast lookups from fullFields
	fullFieldsMap := make(map[string]struct{}, len(fullFields))
	for _, field := range fullFields {
		fullFieldsMap[field] = struct{}{} // Using an empty struct{} as a value, it's memory efficient
	}

	var result []string

	for _, field := range strings.Split(desiredFields, ",") {
		if _, exists := fullFieldsMap[field]; exists {
			result = append(result, field)
		} else if field == "*" {
			for _, fullField := range fullFields {
				result = append(result, fullField)
			}
		}
	}

	return result
}

func validateLimit(limitStr string) uint64 {
	defaultLimit := uint64(10)
	maxLimit := uint64(100)

	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err != nil {
		return defaultLimit
	}

	if limit > maxLimit {
		return maxLimit
	}

	if limit <= 0 {
		return 1
	}

	return limit
}

func validateOffset(pageStr string, limit uint64) uint64 {
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		return uint64(0)
	}

	return (page - 1) * limit
}
