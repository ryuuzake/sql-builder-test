package main

import (
	"context"
	"fmt"
	"net/http"

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
	carColumns := []string{"id", "brand", "model", "year", "state", "color", "fuel_type", "body_type"}
	filterCarColumns := carColumns[1:]

	carsSQL := sq.Select(carColumns...).
		PlaceholderFormat(sq.Dollar).
		From("cars")

	for _, col := range filterCarColumns {
		queryParam := c.QueryParam(col)
		if queryParam != "" {
			carsSQL = carsSQL.Where(sq.Eq{col: queryParam})
		}
	}

	carsSQL = carsSQL.Limit(10)

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
		car, err := scanRowToMap(rows, carColumns)
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
