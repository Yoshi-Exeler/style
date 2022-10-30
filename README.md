# style
a programming style oriented towards real world efficiency

Yoshi Exeler, October 2022

[TOC]

# 0. **Prelude**

## 0.1 **Motivation**
With large parts of the industry fully embracing object oriented design and an almost cult-like following of developers around it, have you ever stopped and taken a step back from the code you were writing, just to realize that you are not actually doing anything productive?

Have you seen codebases getting lost in a sea of bad abstractions, with more boiler plate code than actual business logic?

Have you ever found a `PropertyAuthenticationContainerHelperAttribute` class in a production codebase and wanted to slam your head into your keyboard?

If you answered yes to any of those questions and you have an open mind to alternative approaches, the following paradigm, design philosophy and patterns may help you get back to being productive.

## 0.2 **Programming languages**
The following patterns and design philosophy were developed while writing Go. Depending on the similarity of your language to golang, you may have to adapt some of the more concrete examples into your own language's context. Naturally, some languages will be better suited to this style.

I would strongly recommend you give Go a try. Unlike many other programming languages it is very easy to learn, since it is a small language by design. If you have previous programming experience you should be able to learn Go in less than two Weeks to a point were you are able to write production ready software.

That said, i understand that most programmers don't really have free choice of language, since your company will usually not want to just change languages at the request of a random developer. I hope you will still find some value in the ideas discussed here.

# 0.3 **Base assumptions**
I operate under the following base assumptions, should you disagree with any of these, this style will most likely not be the right thing for you.

## 0.3.1 **Pre-Planning to the implementational level is a bad idea**
While you can plan the structure of an application, define its API's and components, i do not believe you can plan out an entire application (that is non-trivial) down to the implementational level, nor should you. Any application that is of any significance will be so large that the probability of the plan not being correct increases drastically.

Why would the plan not line up with reality?
- Implied functionality was missed during planning (i.e. requirement 'a' says that we need to be able to do 'b' but an implementational detail of 'b' was forgotten)
- A planned implementation is flawed and cannot be implemented
- A requirement of the customer was misunderstood
- The customer forgot a requirement

You might argue that the last two points are not a flaw of the pre-planning approach, since the requirements as they were understood were planned for correctly. While technically correct, you still have to re-plan those sections of the code base, which means that the planning work for the features in question was wasted. 

Therefore, i advocate for an agile, iterative and communication driven development style.

##  0.3.2 **Dynamically typed languages are simply flawed**
Since i believe that static code analysis, IDE tooling and the depth of understanding that the compiler has of your code are vital to a productive programming environment, i am of the opinion that static typing is strictly superior to dynamic typing.
The better the compiler and all your tools can understand your code, the better they can help you in improving it.

Equally important is the fact that static typing forces the programmer to think more deeply about the code they are writing and express their intent more explicitly which leads to code that is more readable and therefore easier to maintain.

# 1. **Core philosophy**

## **1.1 Keep it simple**
For most developers, this guideline is one of the first ones they come into contact with, yet somehow most developers just cannot seem to follow it. I am certainly did my fair share of premature abstraction and unnecessary modularization. To properly follow this rule, one must accept writing boring code.

**How to keep it simple:**
- Do not modularize until you need to
- Do not be afraid to solve concrete problems without abstraction
- Do not optimize things that don't need to be optimized
- If you need to optimize, do algorithmic optimization before implementational optimization

## **1.2 Move fast**
To build software quickly, while staying agile and still ensuring a high code quality, iteration speed is key. Make sure that the time it takes between a developer initiating a build process and the developer either getting feedback or access to their deployed iteration is as small as possible. If programmers need to wait a long time to test their implementations, debugging and testing become a nightmare, which inevitably leads to less testing being done and developer productivity decreasing.

The second component of moving fast is having an environment in which you do not have to be afraid of breaking things. Make sure that your testing/staging/dev environment is setup in such a way that developers do not have to consider damaging anything relevant. Also ensure that your git strategy is such that during a development phase, the developers can move at their own pace and do not have to think about interfering with their colleagues. I would personally recommend either using feature branches or personal branches.

## **1.3 Embrace Static Analysis**
Modern static analysis tools can guarantee certain properties of code, while simultaneously educating developers on how to write better code. Additionally, static analysis can eliminate whole classes of bugs from your code, that you would otherwise have to deal with.

The quality of the static analysis that the tools for your language can provide usually varies greatly by how explicit the language is, which is another huge selling point of statically typed languages that are as explicit as possible.

**Consider the following snippet of javascript:**
```javascript
function logName(v) {
    console.log("Name:" + v.name)
}
```
Since javascript does not require us to specify the type of 'v' this code is unsafe for any types that do not have a name property. Furthermore, when implementing 'logName' our IDE cannot tell us what the properties of 'v' are, meaning that we are just expected to know what they are. This is a fundamental problem, that can only really be fixed by using a programming language with static typing.

**When setting up static analysis for your project, ensure the following things:**
- Developers must be able to run the analysis locally
- Divide your static analysis into two kinds of checks
    - Critical checks are ones that will interrupt the build process if they find any problems
    - Suggestions can be run optionally by developers but are not part of the build process
- nolint/pragma directives must be commented and their comment must contain a concrete plan on how to eventually remove them
- if the number of nolint/pragma directives crosses a certain threshold, developers should drop what they are working on and instead work on removing some directives first

If you are writing Go, look into https://golangci-lint.run/

## **1.4 TODO:Performance by design**

# **2. The Paradigm**
The paradigm defined here aims to be simple and modular with minimal boilerplate code. One of our central goals is to avoid the complex object hierarchies that usually come with object oriented design, since they lead to a number of problems related to reference sharing, state management and initialization.

Note that in the following definitions the words '**trait**' and '**composition**' will be defined in a way not consistent with object oriented design.

## **2.1 Traits**
A trait is a unit of behavior that can be composed with other traits or embedded into data structures.
Unlike inheritance, traits do not implicitly create a hierarchy of data structures, which is something we are actively looking to avoid.

**Defining traits:**
```go
type Sender interface {
	Send(buff []byte) error
}

type Closer interface {
    Close() error
}
```

**Implementing traits:**
```go
type TCPSender struct{}

func (tcps *TCPSender) Send(buff []byte) error { return nil }

type UPDSender struct{}

func (upds *UPDSender) Send(buff []byte) error { return nil }

type QuickCloser struct {}

func (qc *QuickCloser) Close() error { return nil }
```

**Composing traits:**
```go
type ArbitraryConnection struct {
    Sender
    Closer
    ...
}

type GRPCConnection struct {
    TCPSender // here we used struct embedding to ensure this connection uses a tcp sender
    Closer // but we still accept any closer
    ...
}
```

**Accepting trait compositions:**
```go
type Connection interface {
    Sender
    Closer
}

// this method accepts any type that has both the sender and closer trait
func closeConnection(c Connection) {
    c.Close()
}
```

## **2.2 Standalone functions**
Functions may simply exist independently of any other entities. Standalone functions should be pure functions, meaning they should never use global variables or other unsafe shared state. Going forward, standalone functions will be called functions and functions that are associated with an entity will be called methods.

**Consider the following java code example:**
```java
// Why are we doing this? What value does this class provide?
public static class Downloader {
    public static byte[] Download(string url) {
        // do some logic to download the url
    }
}
// even worse, now we need to instantiate a downloader to download something
public class Downloader {
    public byte[] Download(string url) {
        // do some logic to download the url
    }
}
```
A better version of this would be, simply having a function that is not associated to anything else: 
```go
func Download(url string) []byte {
    // do some logic to download the url
}
```
OOP advocates may now argue that instantiating a downloader object makes sense since it allows you to add additional functionality related to downloading later. 

**Consider adding a cache to the downloader: (OOP-Style)**
```java
public class Downloader {
    private Cache cache;

    public Downloader(Cache cache) {
        this.cache = cache;
    }

    public byte[] Download(string url) {
        if this.cache.check(url) {
            return this.cache.yield(url);
        } 
        byte res[] = this.download(url)
        this.cache.populate(url,res);
        return res;
    }

    public byte[] download(string url) {
        // do some logic to download the url
    }
}   
```

**The easiest way to solve this is to implement it explicitly:**
```go
func downloadWithCache(cache Cache) []byte {
    if cache.check(url) {
        return cache.yield(url)
    }
    res := download(url)
    cache.populate(url,res)
    return res
}

func download(url string) []byte {
    // do some logic to download the url
}
```

**And if for whatever reason you need the cache to be encapsulated away you can just use a closure:**
```go
func downloaderWithCache(cache Cache) func(url string) []byte {
    return func(url string) []byte {
        if cache.check(url) {
            return cache.yield(url)
        }
        res := download(url)
        cache.populate(url,res)
        return res
    }
}
```