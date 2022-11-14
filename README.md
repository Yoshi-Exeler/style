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

## **1.4 Performance by design**
When we design a system, we will compose the system from optimized patterns, such that the initial design of the system already plans for the required performance. By considering performance during planning, we prevent the situation of having a system already in place that we need to optimize after its creation. 

Even though we do this, requirements will change over time, so we need to also build our software such that all systems are reasonably encapsulated in packages/functions/methods and have a well defined public interface. If our encapsulation and public api are well designed, we can easily optimize the unit of code at a later time without breaking anything.

Note that this encapsulation can be achieved without requiring speculative generalization, which is one of the main problems of object oriented design.

**Consider the following function snippet:**
```go
func doSomething() {
    ...
    newArray := []string{}
    for _, elem := oldArray range {
        if elem != elementInQuestion {
            newArray = append(newArray,elem)
        }
    }
    ...
}
```
This snippet is something the we could reasonably extract into its own function `removeElement(arr []string, elem string) []string`, allowing us to easily optimize the implementation of the function later without changing anything dependent upon it.


# **2. The Paradigm**
The paradigm defined here aims to be simple and modular with minimal boilerplate code. One of our central goals is to avoid the complex object hierarchies that usually come with object oriented design, since they lead to a number of problems related to reference sharing, state management and initialization.

Note that in the following definitions the words '**trait**' and '**composition**' will be defined in a way not consistent with object oriented design.

## **2.1 Modelling**
The default way of modelling a system must be an explicit implementation. Explicit meaning that we simply plan the minimally abstracted code required to complete a task. If, after modelling such an implementation, we find that there is a good opportunity to reduce the complexity of the planned code using an abstraction, we implement that abstraction. Before implementing an abstraction, try to stop yourself, take a deep breath, take a step back from the current project and objectively evaluate wether the abstraction actually reduces the complexity of the system.

Most importantly, we do not abstract out of laziness, we only abstract when the abstraction is either the only way to reasonably solve the problem or the abstraction reduces the complexity of the codebase. **We do not abstract to save a few lines of code.**

## **2.2 Interfaces & Interface Composition & Interface Embedding**
In cases where we want to create an entity that implements a composition of abstract behaviors, we use interfaces and interface composition.
Interfaces can be used to abstract a concrete behavior. Literal and abstract behaviors may be composed using both struct embedding and interface embedding. This embedding based composition allows us to implement many patterns that would normally be implemented using inheritance. An important distinction is, that this way of implementing Polymorphic entities does not create a hierarchy of entities.

We will further call this style of interface embedding and composition "Traits" and "Trait Composition".

A small example:

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

## **2.3 Data structure composition**
In cases where we want to build an abstraction that combines multiple entities into one greater entity, object oriented design would use inheritance. However, inheritance creates a strict hierarchy of entities, which is one of the things we explicitly try to avoid with this style. So instead of inheritance we use composition. The biggest difference between composition and inheritance is the relationship they create between the involved data types. If a data type 'A' is composed of the data types 'B' and 'C' then 'A' has all fields and methods of 'B' and 'C', but no knowledge of their existence. Meaning that making 'A' a composition of 'B' and 'C' is the same as manually adding the fields and methods of 'B' and 'C' to 'A' as long as 'A', 'B' and 'C' are concrete types. 

This means that given our goals, composing data types is strictly better than using inheritance.
Interestingly, there is the clean code principle `Favour composition over inheritance` (FCoI). Given that even OOP advocates embrace this principle, i propose that we just get rid of inheritance entirely.

## **2.4 Standalone functions**
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
A better version of this would be, simply having a function that is not associated with anything else: 
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
This style of implementation is strictly superior to the OOP-Style for multiple reasons. Our implementation can be easily unit tested, as it does not depend on the state of any entity that is not passed explicitly as a parameter, this means our implementation is strictly functional. In this specific example, you might argue that the OOP-Style implementation is also easily unit-testable if we just model the object as a parameter of the function. While this may work in a small example like this, when dealing with large(r) objects that have many behaviors and contain a lot of state the approach of modeling the object as a parameter becomes less and less feasible as you need to reason about which state is actually relevant to the function. In our style, this work is trivial since all dependencies are passed as parameters.

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

# **3. State Management**
One of the most important problems in modern software development is state management. In my opinion, this problem is the largest failure of object oriented design. 

## **3.1 Three layer state management**
To solve the problem of state management within an application, i propose a state management pattern that divides our code into 3 layers, of which we really only need to deal with two.
### 3.1.1 The Application Kernel
The first layer is the application kernel layer. On this layer there is only one entity, the application kernel. The application kernel must be a singleton with no stateful dependencies. The purpose of the application kernel is to bootstrap our application and provide the ability to spawn asynchronous processes.  

The application bootstrap should work as follows:
- Collect all parameters (CLI, Config, Environment Variables etc)
- Instantiate State Modules and pass Parameters
- Register State Modules with the Kernel
- Call the Kernel Start Method:
  - Compute module load order based on dependencies
  - For each module:
    - Inject dependencies
    - Call module entrypoint

### 3.1.2 State Modules
The second layer of our application is the state module layer. Our application state is managed on this layer. State Modules are modules that manage a closely related chunk of state, exposing a minimal interface to manipulate the state as safely as possible. One of our central design goals is to minimize the amount of code and the size of the api's of the state module layer.


### 3.1.3 Logic Modules
The third layer of our application is the logic module layer. A logic module is a stateless collection of pure functions. As much of our code as possible should be in logic modules. State Modules may call functions from logic modules. Logic modules may only act on state that is explicitly passed to them as a parameter or call other logic code.