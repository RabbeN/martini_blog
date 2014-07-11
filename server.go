package main

import (
	"database/sql"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
)

type Blog struct {
	Id          int
	Title       string
	Date        string
	Description template.HTML
	Author      int
}

func main() {
	m := martini.Classic()
	// render html templates from directory
	m.Use(render.Renderer())

	db, err := sql.Open("mysql", "root:root@unix(/Applications/MAMP/tmp/mysql/mysql.sock)/martini_blog")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	m.Get("/", func(r render.Render) {
		var blogs []Blog

		rows, err := db.Query("select * from blogs")
		if err != nil {
			panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		}
		defer rows.Close()

		for rows.Next() {
			var blog Blog
			var description string
			err := rows.Scan(&blog.Id, &blog.Title, &blog.Date, &description, &blog.Author)
			if err != nil {
				panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
			}
			blog.Description = template.HTML(description)
			blogs = append(blogs, blog)

		}
		err = rows.Err()
		if err != nil {
			panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		}

		r.HTML(200, "home", blogs)
	})

	m.Get("/blog/:id", func(params martini.Params, r render.Render) {
		var blog Blog

		rows, err := db.Query("select * from blogs where id = ?", params["id"])
		if err != nil {
			panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		}
		defer rows.Close()

		for rows.Next() {
			var description string
			err := rows.Scan(&blog.Id, &blog.Title, &blog.Date, &description, &blog.Author)
			if err != nil {
				panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
			}
			blog.Description = template.HTML(description)
		}
		err = rows.Err()
		if err != nil {
			panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		}

		r.HTML(200, "blog", blog)
	})

	m.Run()
}
