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

