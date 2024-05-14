package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Defines a "model" that we can use to communicate with the
// frontend or the database
type BookStore struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	BookName   string
	BookAuthor string
	BookISBN   string
	BookPages  int
	BookYear   int
}

// Wraps the "Template" struct to associate a necessary method
// to determine the rendering procedure
type Template struct {
	tmpl *template.Template
}

// Preload the available templates for the view folder.
// This builds a local "database" of all available "blocks"
// to render upon request, i.e., replace the respective
// variable or expression.
// For more on templating, visit https://jinja.palletsprojects.com/en/3.0.x/templates/
// to get to know more about templating
// You can also read Golang's documentation on their templating
// https://pkg.go.dev/text/template
func loadTemplates() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

// Method definition of the required "Render" to be passed for the Rendering
// engine.
// Contraire to method declaration, such syntax defines methods for a given
// struct. "Interfaces" and "structs" can have methods associated with it.
// The difference lies that interfaces declare methods whether struct only
// implement them, i.e., only define them. Such differentiation is important
// for a compiler to ensure types provide implementations of such methods.
func (t *Template) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

// Here we make sure the connection to the database is correct and initial
// configurations exists. Otherwise, we create the proper database and collection
// we will store the data.
// To ensure correct management of the collection, we create a return a
// reference to the collection to always be used. Make sure if you create other
// files, that you pass the proper value to ensure communication with the
// database
// More on what bson means: https://www.mongodb.com/docs/drivers/go/current/fundamentals/bson/
func prepareDatabase(client *mongo.Client, dbName string, collecName string) (*mongo.Collection, error) {
	db := client.Database(dbName)

	names, err := db.ListCollectionNames(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	if !slices.Contains(names, collecName) {
		cmd := bson.D{{"create", collecName}}
		var result bson.M
		if err = db.RunCommand(context.TODO(), cmd).Decode(&result); err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	coll := db.Collection(collecName)
	return coll, nil
}

// Here we prepare some fictional data and we insert it into the database
// the first time we connect to it. Otherwise, we check if it already exists.
func prepareData(client *mongo.Client, coll *mongo.Collection) {
	startData := []BookStore{
		{
			BookName:   "The Vortex",
			BookAuthor: "José Eustasio Rivera",
			BookISBN:   "958-30-0804-4",
			BookPages:  292,
			BookYear:   1924,
		},
		{
			BookName:   "Frankenstein",
			BookAuthor: "Mary Shelley",
			BookISBN:   "978-3-649-64609-9",
			BookPages:  280,
			BookYear:   1818,
		},
		{
			BookName:   "The Black Cat",
			BookAuthor: "Edgar Allan Poe",
			BookISBN:   "978-3-99168-238-7",
			BookPages:  280,
			BookYear:   1843,
		},
	}

	// This syntax helps us iterate over arrays. It behaves similar to Python
	// However, range always returns a tuple: (idx, elem). You can ignore the idx
	// by using _.
	// In the topic of function returns: sadly, there is no standard on return types from function. Most functions
	// return a tuple with (res, err), but this is not granted. Some functions
	// might return a ret value that includes res and the err, others might have
	// an out parameter.
	for _, book := range startData {
		cursor, err := coll.Find(context.TODO(), book)
		var results []BookStore
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		if len(results) > 1 {
			log.Fatal("more records were found")
		} else if len(results) == 0 {
			result, err := coll.InsertOne(context.TODO(), book)
			if err != nil {
				panic(err)
			} else {
				fmt.Printf("%+v\n", result)
			}

		} else {
			for _, res := range results {
				cursor.Decode(&res)
				fmt.Printf("%+v\n", res)
			}
		}
	}
}

// Generic method to perform "SELECT * FROM BOOKS" (if this was SQL, which
// it is not :D ), and then we convert it into an array of map. In Golang, you
// define a map by writing map[<key type>]<value type>{<key>:<value>}.
// interface{} is a special type in Golang, basically a wildcard...
func findAllBooks(coll *mongo.Collection) []map[string]interface{} {
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	var results []BookStore
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	var ret []map[string]interface{}
	for _, res := range results {
		ret = append(ret, map[string]interface{}{
			"ID":         res.ID.Hex(),
			"BookName":   res.BookName,
			"BookAuthor": res.BookAuthor,
			"BookISBN":   res.BookISBN,
			"BookPages":  res.BookPages,
		})
	}

	return ret
}

func findAllAuthors(coll *mongo.Collection) []string {
    cursor, err := coll.Find(context.TODO(), bson.D{{}})
    var results []BookStore
    if err = cursor.All(context.TODO(), &results); err != nil {
        panic(err)
    }

    var authors []string
    for _, res := range results {
        authors = append(authors, res.BookAuthor)
    }

    return authors
}

func findAllYears(coll *mongo.Collection) []int {
    cursor, err := coll.Find(context.TODO(), bson.D{{}})
    var results []BookStore
    if err = cursor.All(context.TODO(), &results); err != nil {
        panic(err)
    }

    var years []int
    for _, res := range results {
        years = append(years, res.BookYear)
    }

    return years
}

func getAllBooks(coll *mongo.Collection) []map[string]interface{} {
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	var results []BookStore
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	var ret []map[string]interface{}
	for _, res := range results {
		ret = append(ret, map[string]interface{}{
			"id":     res.ID.Hex(),
			"name":   res.BookName,
			"author": res.BookAuthor,
			
			"pages":  res.BookPages,
			"year":   res.BookYear,
			"isbn":   res.BookISBN,
		})
	}

	return ret
}

func main() {
	// Connect to the database. Such defer keywords are used once the local
	// context returns; for this case, the local context is the main function
	// By user defer function, we make sure we don't leave connections
	// dangling despite the program crashing. Isn't this nice? :D
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("DATABASE_URI")
	if len(uri) == 0 {
		fmt.Printf("failure to load env variable\n")
		os.Exit(1)
	}

	// TODO: make sure to pass the proper username, password, and port
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Printf("failed to create client for MongoDB\n")
		os.Exit(1)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Printf("failed to connect to MongoDB, please make sure the database is running\n")
		os.Exit(1)
	}

	// This is another way to specify the call of a function. You can define inline
	// functions (or anonymous functions, similar to the behavior in Python)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// You can use such name for the database and collection, or come up with
	// one by yourself!
	coll, err := prepareDatabase(client, "exercise-2", "information")

	prepareData(client, coll)

	// Here we prepare the server
	e := echo.New()

	// Define our custom renderer
	e.Renderer = loadTemplates()

	// Log the requests. Please have a look at echo's documentation on more
	// middleware
	e.Use(middleware.Logger())

	e.Static("/css", "css")

	// Endpoint definition. Here, we divided into two groups: top-level routes
	// starting with /, which usually serve webpages. For our RESTful endpoints,
	// we prefix the route with /api to indicate more information or resources
	// are available under such route.
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})

	e.GET("/books", func(c echo.Context) error {
		books := findAllBooks(coll)
		return c.Render(200, "book-table", books)
	})

	e.GET("/authors", func(c echo.Context) error {
		authors := findAllAuthors(coll)
		return c.Render(200, "author-table", authors)
	})

	e.GET("/years", func(c echo.Context) error {
		years := findAllYears(coll)
		return c.Render(200, "year-table", years)
	})

	e.GET("/search", func(c echo.Context) error {
		return c.Render(200, "search-bar", nil)
	})

	e.GET("/create", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	e.GET("/api/books", func(c echo.Context) error {
		books := getAllBooks(coll)
		return c.JSON(http.StatusOK, books)
	})

	e.POST("/api/books", func(c echo.Context) error {
		var temp struct {
			Name    string `json:"name"`
			Author  string `json:"author"`
			ISBN    string `json:"isbn"`
			Pages   int    `json:"pages"`
			Year    int    `json:"year"`
		}
	
		if err := c.Bind(&temp); err != nil {
			return err
		}

		book := BookStore{
			BookName:   temp.Name,
			BookAuthor: temp.Author,
			BookISBN:   temp.ISBN,
			BookPages:  temp.Pages,
			BookYear:   temp.Year,
		}
		log.Println("Searching for book:", book)

		if book.BookName == "" || book.BookAuthor == "" || book.BookYear == 0 || book.BookPages == 0 {
			return c.NoContent(http.StatusOK)
		}
		// Check for duplicates
		count, err := coll.CountDocuments(c.Request().Context(), bson.M{
			"bookname":   book.BookName,
			"bookauthor": book.BookAuthor,
			"bookyear":   book.BookYear,
			"bookpages":  book.BookPages,
		})
		log.Println("Count:", count)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusOK, err.Error())
		}

		if count > 0 {
			// A matching document was found
			return echo.NewHTTPError(http.StatusOK, "Book already exists.")
		}
	
		// Insert the new book
		_, err = coll.InsertOne(c.Request().Context(), book)
		if err != nil {
			return echo.NewHTTPError(http.StatusOK, err.Error())
		}
	
		return c.NoContent(http.StatusOK)
	})

	e.DELETE("/api/books/:id", func(c echo.Context) error {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusOK, "Invalid ID")
		}
	
		// Delete the book from the collection
		res, err := coll.DeleteOne(c.Request().Context(), bson.M{"_id": objID})
		if err != nil {
			return echo.NewHTTPError(http.StatusOK, err.Error())
		}
	
		// Check if a book was deleted
		if res.DeletedCount == 0 {
			return echo.NewHTTPError(http.StatusOK, "Book not found")
		}
	
		return c.NoContent(http.StatusOK)
	})

	e.PUT("/api/books", func(c echo.Context) error {
		var temp struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Author string `json:"author"`
			ISBN   string `json:"isbn"`
			Pages  int    `json:"pages"`
			Year   int    `json:"year"`
		}
	
		if err := c.Bind(&temp); err != nil {
			return err
		}
	
		if temp.ID == "" {
			return echo.NewHTTPError(http.StatusOK, "Missing or invalid book ID")
		}
		objID, err := primitive.ObjectIDFromHex(temp.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusOK, "Invalid book ID")
		}
		book := BookStore{
			ID:         objID,
			BookName:   temp.Name,
			BookAuthor: temp.Author,
			BookISBN:   temp.ISBN,
			BookPages:  temp.Pages,
			BookYear:   temp.Year,
		}
		// Get the existing book
		var existingBook BookStore
		err = coll.FindOne(c.Request().Context(), bson.M{"_id": book.ID}).Decode(&existingBook)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusOK, err.Error())
		}

		// Check if any entry has changed
		if existingBook.BookName == book.BookName && existingBook.BookAuthor == book.BookAuthor && existingBook.BookISBN == book.BookISBN && existingBook.BookPages == book.BookPages && existingBook.BookYear == book.BookYear {
			// No entry has changed, do nothing
			return c.NoContent(http.StatusOK)
		}
		// Check for duplicates
		count, err := coll.CountDocuments(c.Request().Context(), bson.M{
			"bookname":   book.BookName,
			"bookauthor": book.BookAuthor,
			"bookyear":   book.BookYear,
			"bookpages":  book.BookPages,
		})
		log.Println("Count:", count)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusOK, err.Error())
		}

		if count > 0 {
			// A matching document was found
			return echo.NewHTTPError(http.StatusOK, "Book already exists.")
		}
	
		// update the book
		_, err = coll.ReplaceOne(c.Request().Context(), bson.M{"_id": book.ID}, book)
		if err != nil {
			return echo.NewHTTPError(http.StatusOK, err.Error())
		}
	
		return c.NoContent(http.StatusOK)
	})
	

	e.Logger.Fatal(e.Start(":3030"))
}
