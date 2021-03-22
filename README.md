Imports are just that: they import code and give you access to identifiers such as types, functions, constants, and interfaces. In our case, the code in the main.go code file can now reference the Run function from the search package, thanks to the import on line 08. On lines 04 and 05, we import code from the standard library for the log and os packages.

Here you see the use of the short variable declaration operator (:=). This operator is used to both declare and initialize variables at the same time. The type of each value being returned is used by the compiler to determine the type for each variable, respectively. The short variable declaration operator is just a shortcut to streamline your code and make the code more readable. The variable it declares is no different than any other variable you may declare when using the keyword var.

```
package search

// import from the standard library
import (
     "log"
     "sync"
)

 // A map of registered matchers for searching.
 var matchers = make(map[string]Matcher)
 
 // The compiler will always look for the packages you import at the locations referenced by the  GOROOT and GOPATH environment variables.
 GOROOT="/Users/me/go"
 GOPATH="/Users/me/spaces/go/projects"
 
 // A map of registered matchers for searching.This variable is located outside the scope of any function and so is considered a package-level variable. The variable is declared using the keyword var and is declared as a map of Matcher type values with a key of type string. The declaration for the Matcher type can be found in the match.go code file. This variable declaration also contains an initialization of the variable via the assignment operator and a special built-in function called make.
 var matchers = make(map[string]Matcher)
 
```

In Go, identifiers are either exported or unexported from a package. An exported identifier can be directly accessed by code in other packages when the respective package is imported. These identifiers start with a capital letter. Unexported identifiers start with a lowercase letter and can’t be directly accessed by code in other packages. But just because an identifier is unexported, it doesn’t mean other packages can’t indirectly access these identifiers. As an example, a function can return a value of an unexported type and this value is accessible by any calling function, even if the calling function has been declared in a different package.

In Go, all variables are initialized to their zero value. For numeric types, that value is 0; for strings it’s an empty string; for Booleans it’s false; and for pointers, the zero value is nil. When it comes to reference types, there are underlying data structures that are initialized to their zero values. But variables declared as a reference type set to their zero value will return the value of nil.

To declare a function in Go, use the keyword func followed by the function name, any parameters, and then any return values. 

Though not unique to Go, you can see that our functions can have multiple return values. It’s common to declare functions that return a value and an error value just like the RetrieveFeeds function. If an error occurs, never trust the other values being returned from the function. They should always be ignored, or else you run the risk of the code generating more errors or panics.

Here you see the use of the short variable declaration operator (:=). This operator is used to both declare and initialize variables at the same time. The type of each value being returned is used by the compiler to determine the type for each variable, respectively. The short variable declaration operator is just a shortcut to streamline your code and make the code more readable. The variable it declares is no different than any other variable you may declare when using the keyword var.

On line 20, we use the built-in function make to create an unbuffered channel. We use the short variable declaration operator to declare and initialize the channel variable with the call to make. A good rule of thumb when declaring variables is to use the keyword var when declaring variables that will be initialized to their zero value, and to use the short variable declaration operator when you’re providing extra initialization or making a function call.

Channels are also a reference type in Go like maps and slices, but channels implement a queue of typed values that are used to communicate data between goroutines. Channels provide inherent synchronization mechanisms to make communication safe. 

In Go, once the main function returns, the program terminates. Any goroutines that were launched and are still running at this time will also be terminated by the Go runtime. When you write concurrent programs, it’s best to cleanly terminate any goroutines that were launched prior to letting the main function return. Writing programs that can cleanly start and shut down helps reduce bugs and prevents resources from corruption.

A *goroutine* is a function that’s launched to run independently from other functions in the program. Use the keyword go to launch and schedule goroutines to run concurrently. An *anonymous function* is a function that’s declared without a name. In our for range loop, we launch an anonymous function as a goroutine for each feed. This allows each feed to be processed independently in a concurrent fashion.

Anonymous functions can take parameters, which we declare for this anonymous function. On line 38 we declare the anonymous function to accept a value of type Matcher and the address of a value of type Feed. This means the variable feed is a *pointer variable*. Pointer variables are great for sharing variables between functions. They allow functions to access and change the state of a variable that was declared within the scope of a different function and possibly a different goroutine.

In Go, all variables are passed by value. Since the value of a pointer variable is the address to the memory being pointed to, passing pointer variables between functions is still considered a pass by value. Thanks to closures, the function can access those variables directly without the need to pass them in as parameters. The anonymous function isn’t given a copy of these variables; it has direct access to the same variables declared in the scope of the outer function. This is the reason why we don’t use closures for the matcher and feed variables.

```
// Run performs the search logic.
12 func Run(searchTerm string) {
13     // Retrieve the list of feeds to search through.

14     feeds, err := RetrieveFeeds()
15     if err != nil {
16         log.Fatal(err)
17     }

18
19     // Create a unbuffered channel to receive match results.
20     results := make(chan *Result)
21
22     // Setup a wait group so we can process all the feeds.
23     var waitGroup sync.WaitGroup
24
25     // Set the number of goroutines we need to wait for while
26     // they process the individual feeds.
27     waitGroup.Add(len(feeds))
28

29     // Launch a goroutine for each feed to find the results.
30     for _, feed := range feeds {
31         // Retrieve a matcher for the search.
32         matcher, exists := matchers[feed.Type]
33         if !exists {
34             matcher = matchers["default"]
35         }
36
37         // Launch the goroutine to perform the search.
38         go func(matcher Matcher, feed *Feed) {
39             Match(matcher, feed, searchTerm, results)
40             waitGroup.Done()
41         }(matcher, feed)
42     }
43
44     // Launch a goroutine to monitor when all the work is done.
45     go func() {
46         // Wait for everything to be processed.
47         waitGroup.Wait()
48
49         // Close the channel to signal to the Display
50         // function that we can exit the program.
51         close(results)
52     }()
53
54     // Start displaying results as they are available and
55     // return after the final result is displayed.
56     Display(results)
57 }
```


```
01 package search
02
03 import (
04     "encoding/json"
05     "os"
06 )
07
08 const dataFile = "data/data.json"
```

```
--- data/data.json
[
    {
        "site" : "npr",
        "link" : "http://www.npr.org/rss/rss.php?id=1001",
        "type" : "rss"
    },
    {
        "site" : "cnn",
        "link" : "http://rss.cnn.com/rss/cnn_world.rss",
        "type" : "rss"
    },
    {
        "site" : "foxnews",
        "link" : "http://feeds.foxnews.com/foxnews/world?format=xml",
        "type" : "rss"
    },
    {
        "site" : "nbcnews",
        "link" : "http://feeds.nbcnews.com/feeds/topstories",
        "type" : "rss"
    }
]
```

These documents need to be decoded into a slice of struct types so we can use this data in our program. Let’s look at the struct type that will be used to decode this data file.

```
10 // Feed contains information we need to process a feed.
11 type Feed struct {
12     Name string `json:"site"`

13     URI  string `json:"link"`
14     Type string `json:"type"`
15 }

17 // RetrieveFeeds reads and unmarshals the feed data file.
18 func RetrieveFeeds() ([]*Feed, error) {
19    // Open the file.
20    file, err := os.Open(dataFile)
21    if err != nil {
22        return nil, err
23    }
24
25    // Schedule the file to be closed once
26    // the function returns.
27    defer file.Close()
28
29    // Decode the file into a slice of pointers
30    // to Feed values.
31    var feeds []*Feed
32    err = json.NewDecoder(file).Decode(&feeds)
33
34    // We don't need to check for errors, the caller can do this.
35    return feeds, err
36 }
```

The keyword defer is used to schedule a function call to be executed right after a function returns. It’s our responsibility to close the file once we’re done with it. By using the keyword defer to schedule the call to the close method, we can guarantee that the method will be called. This will happen even if the function panics and terminates unexpectedly. The keyword defer lets us write this statement close to where the opening of the file occurs, which helps with readability and reducing bugs.

The key to making this code work is the ability of this framework code to use an interface type to capture and call into the specific implementation for each matcher value. This allows the code to handle different types of matcher values in a consistent and generic way. 

search/match.go

```
01 package search
02
03 import (
04     "log"
05 )
06
07 // Result contains the result of a search.
08 type Result struct {
09     Field   string
10     Content string
11 }
12

13 // Matcher defines the behavior required by types that want
14 // to implement a new search type.
15 type Matcher interface {
16     Search(feed *Feed, searchTerm string) ([]*Result, error)
17 }
```

search/default.go

```
01 package search
02
03 // defaultMatcher implements the default matcher.
04 type defaultMatcher struct{}
05
06 // init registers the default matcher with the program.
07 func init() {
08     var matcher defaultMatcher
09     Register("default", matcher)
10 }
11
12 // Search implements the behavior for the default matcher.
13 func (m defaultMatcher) Search(feed *Feed, searchTerm string)
                                                   ([]*Result, error) {
14     return nil, nil
15 }
```

An empty struct allocates zero bytes when values of this type are created. They’re great when you need a type but not any state. For the default matcher, we don’t need to maintain any state; we only need to implement the interface.

Unlike when you call methods directly from values and pointers, when you call a method via an interface type value, the rules are different. Methods declared with pointer receivers can only be called by interface type values that contain pointers. Methods declared with value receivers can be called by interface type values that contain both values and pointers.

```
19 // Match is launched as a goroutine for each individual feed to run
20 // searches concurrently.
21 func Match(matcher Matcher, feed *Feed, searchTerm string,
                                              results chan<- *Result) {
22     // Perform the search against the specified matcher.
23     searchResults, err := matcher.Search(feed, searchTerm)
24     if err != nil {
25         log.Println(err)
26         return
27     }
28
29     // Write the results to the channel.
30     for _, result := range searchResults {
31         results <- result
32     }
33 }

35 // Display writes results to the terminal window as they
36 // are received by the individual goroutines.
37 func Display(results chan *Result) {
38     // The channel blocks until a result is written to the channel.
39     // Once the channel is closed the for loop terminates.
40     for result := range results {
41         fmt.Printf("%s:\n%s\n\n", result.Field, result.Content)
42     }
43 }
```

The structure of the RSS matcher is similar to the structure of the default matcher. It’s the implementation of the interface method Search that’s different and in the end gives each matcher its uniqueness.

```
-- Expected RSS feed document
<rss xmlns:npr="http://www.npr.org/rss/" xmlns:nprml="http://api
    <channel>
        <title>News</title>
        <link>...</link>
        <description>...</description>

        <language>en</language>
        <copyright>Copyright 2014 NPR - For Personal Use
        <image>...</image>
        <item>
            <title>
                Putin Says He'll Respect Ukraine Vote But U.S.
            </title>
            <description>
                The White House and State Department have called on the
            </description>
```

Decoding XML is identical to how we decoded JSON in the feed.go code file.

```
-- matchers/rss.go
01 package matchers
02
03 import (
04     "encoding/xml"
05     "errors"
06     "fmt"
07     "log"
08     "net/http"
09     "regexp"
10
11     "github.com/goinaction/code/chapter2/sample/search"
12 )

14 type (
15     // item defines the fields associated with the item tag
16     // in the rss document.
17     item struct {
18         XMLName     xml.Name `xml:"item"`
19         PubDate     string   `xml:"pubDate"`
20         Title       string   `xml:"title"`
21         Description string   `xml:"description"`

22         Link        string   `xml:"link"`
23         GUID        string   `xml:"guid"`
24         GeoRssPoint string   `xml:"georss:point"`
25     }
26
27     // image defines the fields associated with the image tag
28     // in the rss document.
29     image struct {
30         XMLName xml.Name `xml:"image"`
31         URL     string   `xml:"url"`
32         Title   string   `xml:"title"`
33         Link    string   `xml:"link"`
34     }
35
36     // channel defines the fields associated with the channel tag
37     // in the rss document.
38     channel struct {
39         XMLName        xml.Name `xml:"channel"`
40         Title          string   `xml:"title"`
41         Description    string   `xml:"description"`
42         Link           string   `xml:"link"`
43         PubDate        string   `xml:"pubDate"`
44         LastBuildDate  string   `xml:"lastBuildDate"`
45         TTL            string   `xml:"ttl"`
46         Language       string   `xml:"language"`
47         ManagingEditor string   `xml:"managingEditor"`
48         WebMaster      string   `xml:"webMaster"`
49         Image          image    `xml:"image"`
50         Item           []item   `xml:"item"`
51    }
52
53    // rssDocument defines the fields associated with the rss document
54    rssDocument struct {
55         XMLName xml.Name `xml:"rss"`
56         Channel channel  `xml:"channel"`
57    }
58 )

60 // rssMatcher implements the Matcher interface.
61 type rssMatcher struct{}

63 // init registers the matcher with the program.
64 func init() {
65     var matcher rssMatcher
66     search.Register("rss", matcher)
67 }
```

The unexported method retrieve performs the logic for pulling the RSS document from the web for each individual feed link. On line 121 you can see the use of the Get method from the http package.

```
-- matchers/rss.go
114 // retrieve performs a HTTP Get request for the rss feed and decodes
115 func (m rssMatcher) retrieve(feed *search.Feed)
                                                 (*rssDocument, error) {
116     if feed.URI == "" {
117         return nil, errors.New("No rss feed URI provided")
118     }
119
120     // Retrieve the rss feed document from the web.
121     resp, err := http.Get(feed.URI)
122     if err != nil {
123         return nil, err
124     }
125
126     // Close the response once we return from the function.
127     defer resp.Body.Close()
128
129     // Check the status code for a 200 so we know we have received a
130     // proper response.
131     if resp.StatusCode != 200 {
132         return nil, fmt.Errorf("HTTP Response Error %d\n",
                                                        resp.StatusCode)
133     }
134
135     // Decode the rss feed document into our struct type.

136     // We don't need to check for errors, the caller can do this.
137     var document rssDocument
138     err = xml.NewDecoder(resp.Body).Decode(&document)
139     return &document, err
140 }
```


```
matchers/rss.go
 69 // Search looks at the document for the specified search term.
 70 func (m rssMatcher) Search(feed *search.Feed, searchTerm string)
                                            ([]*search.Result, error) {
 71     var results []*search.Result
 72
 73     log.Printf("Search Feed Type[%s] Site[%s] For Uri[%s]\n",
                                        feed.Type, feed.Name, feed.URI)
 74
 75     // Retrieve the data to search.
 76     document, err := m.retrieve(feed)
 77     if err != nil {
 78         return nil, err
 79     }
 80
 81     for _, channelItem := range document.Channel.Item {
 82         // Check the title for the search term.
 83         matched, err := regexp.MatchString(searchTerm,
                                                     channelItem.Title)
 84         if err != nil {
 85             return nil, err
 86         }
 87
 88         // If we found a match save the result.
 89         if matched {
 90            results = append(results, &search.Result{
 91                Field:   "Title",
 92                Content: channelItem.Title,

 93            })
 94         }
 95
 96         // Check the description for the search term.
 97         matched, err = regexp.MatchString(searchTerm,
                                               channelItem.Description)
 98         if err != nil {
 99             return nil, err
100         }
101
102         // If we found a match save the result.
103         if matched {
104             results = append(results, &search.Result{
105                 Field:   "Description",
106                 Content: channelItem.Description,
107             })
108         }
109     }
110
111     return results, nil
112 }
```

#### 2.5. SUMMARY

- Every code file belongs to a package, and that package name should be the same as the folder the code file exists in.
- Go provides several ways to declare and initialize variables. If the value of a variable isn’t explicitly initialized, the compiler will initialize the variable to its zero value.
- Pointers are a way of sharing data across functions and goroutines.
- Concurrency and synchronization are accomplished by launching goroutines and using channels.
- Go provides built-in functions to support using Go’s internal data structures.
- The standard library contains many packages that will let you do some powerful things.
- Interfaces in Go allow you to write generic code and frameworks.

**Packages**

All Go programs are organized into groups of files called *packages*, so that code has the ability to be included into other projects as smaller reusable pieces. Each package can be imported and used individually so that developers can import only the specific functionality that they need. This means that all .go files in a single directory must declare the same package name.

```
net/http/
    cgi/
    cookiejar/
        testdata/
    fcgi/
    httptest/
    httputil/
    pprof/
    testdata/
```

The package name main has special meaning in Go. It designates to the Go command that this package is intended to be compiled into a binary executable. All of the executable programs you build in Go must have a package called main.

When the main package is encountered by the compiler, it must also find a function called main(); otherwise a binary executable won’t be created. The main() function is the entry point for the program, so without one, the program has no starting point. The name of the final binary will take the name of the directory the main package is declared in.

The Go documentation uses the term *command* frequently to refer to an executable program—like a command-line application. Remember that in Go, a command is any executable program, in contrast to a package, which generally means an importable semantic unit of functionality.

**Imports**

 The import statement tells the compiler where to look on disk to find the package you want to import. Packages are found on disk based on their relative path to the directories referenced by the Go environment. Packages in the standard library are found under where Go is installed on your computer. Packages that are created by you or other Go developers live inside the GOPATH, which is your own personal workspace for packages.

When an import path contains a URL, the Go tooling can be used to fetch the package from the DVCS and place the code inside the GOPATH at the location that matches the URL. This fetching is done using the go get command. go get will fetch any specified URL or can be used to fetch the dependencies a package is importing that are go-gettable. Since go get is recursive, it can walk down the source tree for a package and fetch all the dependencies it finds.

The _ (underscore character) is known as the *blank identifier* and has many uses within Go. It’s used when you want to throw away the assignment of a value, including the assignment of an import to its package name, or ignore return values from a function when you’re only interested in the others.

**Init**

Each package has the ability to provide as many init functions as necessary to be invoked at the beginning of execution time. All the init functions that are discovered by the compiler are scheduled to be executed prior to the main function being executed. The init functions are great for setting up packages, initializing variables, or performing any other bootstrapping you may need prior to the program running.

```
// Sample program to show how to show you how to briefly work
// with the sql package.
package main

import (
	"database/sql"

	_ "github.com/goinaction/code/chapter3/dbdriver/postgres"
)

// main is the entry point for the application.
func main() {
	sql.Open("postgres", "mydb")
}
```

Go tools

```go build hello.go
go build hello.go

go clean hello.go

go run wordcount.go

go vet wordcount.go

go fmt wordcount.go

go doc tar
```

Go vet for styling errors:

- Bad parameters in Printf-style function calls
- Method signature errors for common method definitions
- Bad struct tags
- Unkeyed composite literals

Browsing the documentation

```
godoc -http=:6060
```

This command instructs godoc to start a web server on port 6060. If you open your web browser and navigate to http://localhost:6060, you’ll see a web page with documentation for both the Go standard libraries and any Go source that lives in your GOPATH.

Once you start cranking out awesome Go code, you’re probably going to want to share that code with the rest of the Go community. It’s really easy as long as you follow a few simple steps.

**Package should live at the root of the repository**

When you’re using go get, you specify the full path to the package that should be imported. This means that when you create a repository that you intend to share, the package name should be the repository name, and the package’s source should be in the root of the repository’s directory structure.

A common mistake that new Go developers make is to create a code or src directory in their public repository. Doing so will make the package’s public import longer. Instead, just put the package source files at the root of the public repository.

**Packages can be small**

It’s common in Go to see packages that are relatively small by the standards of other programming languages. Don’t be afraid to make a package that has a small API or performs only a single task. That’s normal and expected.

**Run go fmt on the code**

Just like any other open source repository, people will look at your code to gauge the quality of it before they try it out. You need to be running go fmt before checking anything in. It makes your code readable and puts everyone on the same page when reading source code.

**Document the code**

Go developers use godoc to read documentation, and [http://godoc.org](http://godoc.org/) to read documentation for open source packages. If you’ve followed go doc best practices in documenting your code, your packages will appear well documented when viewed locally or online, and people will find it easier to use.

**Dependency**

**Vendoring dependencies**

Community tools such as godep and vendor have solved the dependency problem by using a technique called *vendoring* and import path rewriting. The idea is to copy all the dependencies into a directory inside the project repo, and then rewrite any import paths that reference those dependencies by providing the location inside the project itself.

Before the dependencies were vendored, the import statements used the canonical path for the package. The code was physically located on disk within the scope of GOPATH. After vendoring, import path rewriting became necessary to reference the packages, which are now physically located on disk inside the project itself. You can see these imports are very large and tedious to use.

With vendoring, you have the ability to create reproducible builds, since all the source code required to build the binary is housed inside the single project repo. One other benefit of vendoring and import path rewriting is that the project repo is still go-gettable. When go get is called against the project repo, the tooling can find each package and store the package exactly where it needs to be inside the project itself.

**Gb**

The philosophy behind gb stems from the idea that Go doesn’t have reproducible builds because of the import statement. The import statement drives go get, but import doesn’t contain sufficient information to identify which revision of a package should be fetched any time go get is called. The possibility that go get can fetch a different version of code for any given package at any time makes supporting the Go tooling in any reproducible solution complicated and tedious at best. You saw some of this tediousness with the large import paths when using godep.

Gb doesn’t wrap the Go tooling, nor does it use GOPATH. Gb replaces the Go tooling workspace metaphor with a project-based approach. This has natively allowed vendoring without the need for rewriting import paths, which is mandated by go get and a GOPATH workspace.

Gb projects differentiate between the code you write and the code your code depends on. The code your code depends on is called *vendored code*. A gb project makes a clear distinction between your code and vendored code.

One of the best things about gb is that there’s no need for import path rewriting. Look at the import statements that are declared inside of main.go—nothing needs to change to reference the vendored dependencies.

The gb tool will look inside the ```$PROJECT/vendor/src/``` directory for these imports if they can’t be located inside the $PROJECT/src/ directory first. The entire source code for the project is located within a single repo and directory on disk, split between the src/ and vendor/src/ subdirectories. This, in conjunction with no need to rewrite import paths and the freedom to place your project anywhere you wish on disk, makes gb a popular tool in the community to develop projects that require reproducible builds.

One thing to note: a gb project is not compatible with the Go tooling, including go get. Since there’s no need for GOPATH, and the Go tooling doesn’t understand the structure of a gb project, it can’t be used to build, test, or get. Building and testing a gb project requires navigating to the $PROJECT directory and using the gb tool.

Many of the same features that are supported by the Go tooling are supported in gb. Gb also has a plugin system to allow the community to extend support. One such plugin is called vendor, which provides conveniences to manage the dependencies in the ```$PROJECT/vendor/src/ directory```, something the Go tooling does not have today. To learn more about gb, check out the website: getgb.io.

#### 3.8. SUMMARY

- Packages are the basic unit of code organization in Go.
- Your GOPATH determines on disk where Go source code is saved, compiled, and installed.
- You can set your GOPATH for each different project, keeping all of your source and dependencies separate.
- The go tool is your best friend when working from the command line.
- You can use packages created by other people by using go get to fetch and install them in your GOPATH.
- It’s easy to create packages for others to use if you host them on a public source code repository and follow a few simple rules.
- Go was designed with code sharing as a central driving feature of the language.
- It’s recommended that you use vendoring to manage dependencies.
- There are several community-developed tools for dependency management such as godep, vendor, and gb.

**Array intervals and Fundamentals**

An array in Go is a fixed-length data type that contains a contiguous block of elements of the same type. This could be a built-in type such as integers and strings, or it can be a struct type. Arrays are valuable data structures because the memory is allocated sequentially. Having memory in a contiguous form can help to keep the memory you use stay loaded within CPU caches longer. Using index arithmetic, you can iterate through all the elements of an array quickly. The type information for the array provides the distance in memory you have to move to find each element. Since each element is of the same type and follows each other sequentially, moving through the array is consistent and fast.

Once an array is declared, neither the type of data being stored nor its length can be changed. If you need more elements, you need to create a new array with the length needed and then copy the values from one array to the other.

If the length is given as ..., Go will identify the length of the array based on the number of elements that are initialized.

You can have an array of pointers, you use the * operator to access the value that each element pointer points to.

An array is a value in Go. This means you can use it in an assignment operation. The variable name denotes the entire array and, therefore, an array can be assigned to other arrays of the same type.

```
-- 4.1. Declaring an array set to its zero value
// Declare an integer array of five elements.
var array [5]int

-- 4.2. Declaring an array using an array literal
// Declare an integer array of five elements.
// Initialize each element with a specific value.
array := [5]int{10, 20, 30, 40, 50}


-- 4.3. Declaring an array with Go calculating size
// Declare an integer array.
// Initialize each element with a specific value.
// Capacity is determined based on the number of values initialized.
array := [...]int{10, 20, 30, 40, 50}

- 4.4. Declaring an array initializing specific elements
// Declare an integer array of five elements.
// Initialize index 1 and 2 with specific values.
// The rest of the elements contain their zero value.
array := [5]int{1: 10, 2: 20}

-- 4.5. Accessing array elements
// Declare an integer array of five elements.
// Initialize each element with a specific value.
array := [5]int{10, 20, 30, 40, 50}


// Change the value at index 2.
array[2] = 35

-- 4.6. Accessing array pointer elements
// Declare an integer pointer array of five elements.
// Initialize index 0 and 1 of the array with integer pointers.
array := [5]*int{0: new(int), 1: new(int)}

// Assign values to index 0 and 1.
*array[0] = 10
*array[1] = 20

--  4.7. Assigning one array to another of the same type
// Declare a string array of five elements.
var array1 [5]string

-- 4.8. Compiler error assigning arrays of different types
// Declare a string array of four elements.
var array1 [4]string

// Declare a second string array of five elements.
// Initialize the array with colors.
array2 := [5]string{"Red", "Blue", "Green", "Yellow", "Pink"}

// Copy the values from array2 into array1.
array1 = array2

Compiler Error:
cannot use array2 (type [5]string) as type [4]string in assignment

-- 4.9. Assigning one array of pointers to another
// Declare a string pointer array of three elements.
var array1 [3]*string

// Declare a second string pointer array of three elements.
// Initialize the array with string pointers.
array2 := [3]*string{new(string), new(string), new(string)}

// Add colors to each element
*array2[0] = "Red"
*array2[1] = "Blue"
*array2[2] = "Green"

// Copy the values from array2 into array1.
array1 = array2

-- 4.10. Declaring two-dimensional arrays
// Declare a two dimensional integer array of four elements
// by two elements.
var array [4][2]int

// Use an array literal to declare and initialize a two
// dimensional integer array.
array := [4][2]int{{10, 11}, {20, 21}, {30, 31}, {40, 41}}

// Declare and initialize index 1 and 3 of the outer array.
array := [4][2]int{1: {20, 21}, 3: {40, 41}}

// Declare and initialize individual elements of the outer
// and inner array.
array := [4][2]int{1: {0: 20}, 3: {1: 41}}

-- 4.11. Accessing elements of a two-dimensional array
// Declare a two dimensional integer array of two elements.
var array [2][2]int

// Set integer values to each individual element.
array[0][0] = 10
array[0][1] = 20
array[1][0] = 30
array[1][1] = 40

-- 4.12. Assigning multidimensional arrays of the same type
// Declare two different two dimensional integer arrays.
var array1 [2][2]int
var array2 [2][2]int

// Add integer values to each individual element.
array2[0][0] = 10
array2[0][1] = 20
array2[1][0] = 30
array2[1][1] = 40

// Copy the values from array2 into array1.
array1 = array2

-- 4.13. Assigning multidimensional arrays by index
// Copy index 1 of array1 into a new array of the same type.
var array3 [2]int = array1[1]

// Copy the integer found in index 1 of the outer array
// and index 0 of the interior array into a new variable of
// type integer.
var value int = array1[1][0]

-- 4.14. Passing a large array by value between functions
// Declare an array of 8 megabytes.
var array [1e6]int

// Pass the array to the function foo.
foo(array)

// Function foo accepts an array of one million integers.
func foo(array [1e6]int) {
    ...
}

-- 4.15. Passing a large array by pointer between functions
// Allocate an array of 8 megabytes.
var array [1e6]int

// Pass the address of the array to the function foo.
foo(&array)

// Function foo accepts a pointer to an array of one million integers.
func foo(array *[1e6]int) {
    ...
}
```

**SLICE INTERNALS AND FUNDAMENTALS**

A *slice* is a data structure that provides a way for you to work with and manage collections of data. Slices are built around the concept of dynamic arrays that can grow and shrink as you see fit. They’re flexible in terms of growth because they have their own built-in function called append, which can grow a slice quickly with efficiency. You can also reduce the size of a slice by slicing out a part of the underlying memory. Slices give you all the benefits of indexing, iteration, and garbage collection optimizations because the underlying memory is allocated in contiguous blocks.

**Internals**

```
-- 4.16. Declaring a slice of strings by length
// Create a slice of strings.
// Contains a length and capacity of 5 elements.
slice := make([]string, 5)

-- 4.17. Declaring a slice of integers by length and capacity
// Create a slice of integers.
// Contains a length of 3 and has a capacity of 5 elements.
slice := make([]int, 3, 5)

-- 4.18. Compiler error setting capacity less than length
// Create a slice of integers.
// Make the length larger than the capacity.
slice := make([]int, 5, 3)

Compiler Error:
len larger than cap in make([]int)

-- 4.19. Declaring a slice with a slice literal
// Create a slice of strings.
// Contains a length and capacity of 5 elements.
slice := []string{"Red", "Blue", "Green", "Yellow", "Pink"}

// Create a slice of integers.
// Contains a length and capacity of 3 elements.
slice := []int{10, 20, 30}

-- 4.21 Declaration differences between arrays and slices
// Create an array of three integers.
array := [3]int{10, 20, 30}

// Create a slice of integers with a length and capacity of three.
slice := []int{10, 20, 30}

-- 4.22. Declaring a nil slice
// Create a nil slice of integers.
var slice []int

-- 4.23. Declaring an empty slice
// Use make to create an empty slice of integers.
slice := make([]int, 0)

// Use a slice literal to create an empty slice of integers.
slice := []int{}
```

The three fields are a pointer to the underlying array, the length or the number of elements the slice has access to, and the capacity or the number of elements the slice has available for growth. The difference between length and capacity will make more sense in a bit.

When you specify the length and capacity separately, you can create a slice with available capacity in the underlying array that you don’t have access to initially.

 if you specify a value inside the [ ] operator, you’re creating an array. If you don’t specify a value, you’re creating a slice.

Sometimes in your programs you may need to declare a nil slice. A nil slice is created by declaring a slice without any initialization. A nil slice is the most common way you create slices in Go. They can be used with many of the standard library and built-in functions that work with slices. 

An empty slice contains a zero-element underlying array that allocates no storage. Empty slices are useful when you want to represent an empty collection, such as when a database query returns zero results 

```
-- 4.24. Declaring an array using an array literal
// Create a slice of integers.
// Contains a length and capacity of 5 elements.
slice := []int{10, 20, 30, 40, 50}

// Change the value of index 1.
slice[1] = 25

-- 4.25. Taking the slice of a slice
// Create a slice of integers.
// Contains a length and capacity of 5 elements.
slice := []int{10, 20, 30, 40, 50}


// Create a new slice.
// Contains a length of 2 and capacity of 4 elements.
newSlice := slice[1:3]

-- 4.26. How length and capacity are calculated
For slice[i:j] with an underlying array of capacity k

Length:   j - i
Capacity: k - i

-- 4.28. Potential consequence of making changes to a slice
// Create a slice of integers.
// Contains a length and capacity of 5 elements.
slice := []int{10, 20, 30, 40, 50}

// Create a new slice.
// Contains a length of 2 and capacity of 4 elements.
newSlice := slice[1:3]

// Change index 1 of newSlice.
// Change index 2 of the original slice.
newSlice[1] = 35

-- 4.29. Runtime error showing index out of range
// Create a slice of integers.
// Contains a length and capacity of 5 elements.
slice := []int{10, 20, 30, 40, 50}

// Create a new slice.
// Contains a length of 2 and capacity of 4 elements.
newSlice := slice[1:3]

// Change index 3 of newSlice.
// This element does not exist for newSlice.
newSlice[3] = 45

Runtime Exception:
panic: runtime error: index out of range
```

Slices are called such because you can slice a portion of the underlying array to create a new slice.

Having capacity is great, but useless if you can’t incorporate it into your slice’s length. Luckily, Go makes this easy when you use the built-in function append. 

One of the advantages of using a slice over using an array is that you can grow the capacity of your slice as needed. Go takes care of all the operational details when you use the built-in function append.

To use append, you need a source slice and a value that is to be appended. When your append call returns, it provides you a new slice with the changes. The append function will always increase the length of the new slice. The capacity, on the other hand, may or may not be affected, depending on the available capacity of the source slice.

When there’s no available capacity in the underlying array for a slice, the append function will create a new underlying array, copy the existing values that are being referenced, and assign the new value.

The append operation is clever when growing the capacity of the underlying array. Capacity is always doubled when the existing capacity of the slice is under 1,000 elements. Once the number of elements goes over 1,000, the capacity is grown by a factor of 1.25, or 25%. This growth algorithm may change in the language over time.

The built-in function append will use any available capacity first. Once that capacity is reached, it will allocate a new underlying array. It’s easy to forget which slices are sharing the same underlying array. When this happens, making changes to a slice can result in random and odd-looking bugs. Suddenly changes appear on multiple slices out of nowhere.

By having the option to set the capacity of a new slice to be the same as the length, you can force the first append operation to detach the new slice from the underlying array. Detaching the new slice from its original source array makes it safe to change.

With the new slice now having its own underlying array, we’ve avoided potential problems. We can now continue to append fruit to our new slice without worrying if we’re changing fruit to other slices inappropriately. Also, allocating the new underlying array for the slice was easy and clean.

There are two special built-in functions called len and cap that work with arrays, slices, and channels. For slices, the len function returns the length of the slice, and the cap function returns the capacity. 

```
-- 4.30. Using append to add an element to a slice
// Create a slice of integers.
// Contains a length and capacity of 5 elements.
slice := []int{10, 20, 30, 40, 50}

// Create a new slice.
// Contains a length of 2 and capacity of 4 elements.
newSlice := slice[1:3]

// Allocate a new element from capacity.
// Assign the value of 60 to the new element.
newSlice = append(newSlice, 60)

-- 4.31. Using append to increase the length and capacity of a slice
// Create a slice of integers.
// Contains a length and capacity of 4 elements.
slice := []int{10, 20, 30, 40}

// Append a new value to the slice.
// Assign the value of 50 to the new element.
newSlice := append(slice, 50)

-- 4.32. Declaring a slice of string using a slice literal
// Create a slice of strings.
// Contains a length and capacity of 5 elements.
source := []string{"Apple", "Orange", "Plum", "Banana", "Grape"}

-- 4.33. Performing a three-index slice
// Slice the third element and restrict the capacity.
// Contains a length of 1 element and capacity of 2 elements.
slice := source[2:3:4]

-- 4.36. Benefits of setting length and capacity to be the same
// Create a slice of strings.
// Contains a length and capacity of 5 elements.
source := []string{"Apple", "Orange", "Plum", "Banana", "Grape"}

// Slice the third element and restrict the capacity.
// Contains a length and capacity of 1 element.
slice := source[2:3:3]

// Append a new string to the slice.
slice = append(slice, "Kiwi")

-- 4.39. range provides a copy of each element
// Create a slice of integers.
// Contains a length and capacity of 4 elements.
slice := []int{10, 20, 30, 40}

// Iterate over each element and display the value and addresses.
for index, value := range slice {
   fmt.Printf("Value: %d  Value-Addr: %X  ElemAddr: %X\n",
       value, &value, &slice[index])
}

-- 4.40. Using the blank identifier to ignore the index value
// Create a slice of integers.
// Contains a length and capacity of 4 elements.
slice := []int{10, 20, 30, 40}

// Iterate over each element and display each value.
for _, value := range slice {
    fmt.Printf("Value: %d\n", value)
}

-- 4.41. Iterating over a slice using a traditional for loop
// Create a slice of integers.
// Contains a length and capacity of 4 elements.
slice := []int{10, 20, 30, 40}

// Iterate over each element starting at element 3.
for index := 2; index < len(slice); index++ {
    fmt.Printf("Index: %d  Value: %d\n", index, slice[index])
}
```
