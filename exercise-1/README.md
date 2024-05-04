## Exercise 1 for Cloud Computing ##

### Summer Semester 2024 ###

#### More enquiries, please contact me ####

Dear students, 

Let's embark in this journey exploring the development of a simple, yet very dynamic webpage. During this time, we will combine cloud standard technologies (Golang) and emerging ones (HTMX) to discover what it takes to deploy an application in the Cloud. From a simple binary running to using FaaS for efficiency and cost optimization.

In this first week, we will explore developing certain views using templating and content swap in the browser. All the heavy-lifting is done by HTMX, you just need to adjust the contents and attributes of the *HTML tags* to signal how it should react.

In this code, there is already quite a bit of examples on how you can do it: from the usage of templates, rendering upon request, to interacting with the database to maintain your data.

In the topic of database: there is a bunch of different databases, as there is *database engines*. For simplicity, we will use a NoSQL database called MongoDB. Why NoSQL? Well, the name implies Non-relational database, i.e., here we don't care that much about those peski relations between tables, and foreing keys, and all that technical jargon. Basically, MongoDB works like a big storage for "JSON" documents (or more known as Key-value stores a.k.a KV-stores). That means, for this exercise, we want something simple yet standard in Cloud environments, just to get dabble with them and explore what we can do.

### What is the task? ###

As you have seen from the exercise slides or recording, we will be exploring developing a web-front for a Book Store. For this task, you will need to complete:

1. Finish the implementation for the "Authors" view by: implementing the HTML template, and rendering the content upon a request to `/authors`
2. Finish the implementation for the "Years" view by: implementing the HTML template, and rendering the content upon a request to `/years`
3. Implement the methods for:

    3.1. `GET`. The request path should be `/api/books`, and it should return an array of objects, in the following form:

        response = {
                id: "asd34343",
                name: "The book name",
                author: "The book author",
                pages: 1000,
                year: 1900,
                isbn: "900-01-45-2222" // this field is optional
        }

    3.2. `POST`. The request path should be `/api/books`, and it should return the status code 200 upon **correct** completion. The body of the request looks as follows: 

        request.body = {
                name: "The book name",
                author: "The book author",
                pages: 1000,
                year: 1900,
                isbn: "900-01-45-2222" // this field is optional
        }

    3.3 `UPDATE`. The request path should be `/api/books`, and it should return the status code 200 upon **correct** completion. The body of the request looks as follows:

        request.body = {
                id: "asd34343",
                name: "The book name",
                author: "The book author",
                pages: 1000,
                year: 1900,
                isbn: "900-01-45-2222" // this field is optional
        }

    3.4 `DELETE`. The request path should be `/api/books/:id`, and it should return the status code 200 upon **correct** deletion of the respective book. In this context, `:id` is known as a path parameter, and common HTTP server frameworks (like the one we are using), supports parsing such parameter to the point you can easily access it. The value for `:id` is the key `id` from previous responses.

    3.5 We will make **six** tests to these endpoints with random (but stable data) to make sure the workflow is correct. If everything is correct, together with the rendering functionality, you will achieve 100 points. Remember that you need **70 %** to pass each assignment.


### What do you need for this assignment? ###

Basically, only one thing: **Golang 1.22**. You can download and install it following [these](https://go.dev/doc/install) instructions.

Once downloaded and installed, under the exercise folder, run the following commands:

> go mod tidy // it will automatically download all dependencies specified in go.mod and go.sum

> go run main.go // it will launch the server and let you access it via localhost:3030

To build your binary, you can perform the following command:

> go build -o <out_filename>

The other component you need to run your exercise is a database. Since we are using MongoDB, you can installing following the instructions [here](https://www.mongodb.com/docs/v7.0/administration/install-community/). I recommend you use MongoDB CE v.7. Moreover, you will also have to change the MongoDB host inside [main.go](cmd/main.go#L184). Remember that you must also specify an username and password when installing MongoDB. In my case, I chose `mongodb` as user, and `testmongo` as password. The port in the [URI](https://en.wikipedia.org/wiki/Uniform_Resource_Identifier) must be also replace to match your system.

Without further ado,

#### Happy Coding! ####


Current Mongosh Log ID: 6634026479fa519d6d2202d7
Connecting to:          mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.2.5
Using MongoDB:          7.0.9
Using Mongosh:          2.2.5

For mongosh info see: https://docs.mongodb.com/mongodb-shell/


To help improve our products, anonymous usage data is collected and sent to MongoDB periodically (https://www.mongodb.com/legal/privacy-policy).      
You can opt-out by running the disableTelemetry() command.

------
   The server generated these startup warnings when booting
   2024-05-02T23:14:24.625+02:00: Using the XFS filesystem is strongly recommended with the WiredTiger storage engine. See http://dochub.mongodb.org/core/prodnotes-filesystem
   2024-05-02T23:14:24.820+02:00: Access control is not enabled for the database. Read and write access to data and configuration is unrestricted     
   2024-05-02T23:14:24.820+02:00: /sys/kernel/mm/transparent_hugepage/enabled is 'always'. We suggest setting it to 'never' in this binary version    
   2024-05-02T23:14:24.820+02:00: vm.max_map_count is too low
------