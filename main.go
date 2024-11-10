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
	brand := c.QueryParam("brand")
	model := c.QueryParam("model")
	year := c.QueryParam("year")
	state := c.QueryParam("state")
	color := c.QueryParam("color")
	fuel_type := c.QueryParam("fuel_type")
	body_type := c.QueryParam("body_type")

	var cars []Car

	carsSQL := sq.Select("id", "brand", "model", "year", "state", "color", "fuel_type", "body_type").
		PlaceholderFormat(sq.Dollar).
		From("cars")

	if brand != "" {
		carsSQL = carsSQL.Where(sq.Eq{"brand": brand})
	}

	if model != "" {
		carsSQL = carsSQL.Where(sq.Eq{"model": model})
	}

	if year != "" {
		carsSQL = carsSQL.Where(sq.Eq{"year": year})
	}

	if state != "" {
		carsSQL = carsSQL.Where(sq.Eq{"state": state})
	}

	if color != "" {
		carsSQL = carsSQL.Where(sq.Eq{"color": color})
	}

	if fuel_type != "" {
		carsSQL = carsSQL.Where(sq.Eq{"fuel_type": fuel_type})
	}

	if body_type != "" {
		carsSQL = carsSQL.Where(sq.Eq{"body_type": body_type})
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

	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.ID, &car.Brand, &car.Model, &car.Year, &car.State, &car.Color, &car.FuelType, &car.BodyType); err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to scan car"})
		}

		cars = append(cars, car)
	}

	return c.JSON(http.StatusOK, cars)
}
