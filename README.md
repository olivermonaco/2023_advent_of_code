# Retro

## What went well

### Bash script templating the challenges
- Aside from being satisfying to automate, Bash script templating helped me get into challenges quickly, without having to set up much boilerplate
- The setup didn't take too long, and it was probably worth it from a time perspective (see [in depth research](https://xkcd.com/1205) on the subject). Time costs aside, the real gain for me was alleviating some of the mental startup cost of beginning a new coding challenge each time

### Forgiving myself for getting behind
- I can be a bit completist, and the holiday season (and life generally) is busy. I realized that I might not be able to keep to a satisfying pace, and have a life. I spent as much time as I needed to without going overboard on trying to complete every day in the calendar

### Logging and CSV creation
- multiple outputs (to std out and a file) for logs. Helped in for slightly more complex cases
- using the context and subloggers instead of a global logger. My logging tool of choice, [zerolog](https://pkg.go.dev/github.com/rs/zerolog), is an additive logger. Associating the logger with the context and instantiating subloggers down the call stack allowed me to scope  logging information to the current case. This helps to avoid the confusion of adding details from other cases
- csv file creation for when I wanted greater detail, and the logging file format wasn't clear enough

### Panicking, and other quick and dirty approaches
- I'm not talking about panicking in the anxiety attack sense, but using [panics](https://gobyexample.com/panic) in Go. And, overall I'm talking about deciding to code in a quick and dirty way. 
- When I write code for work my goal is to write well documented, tested, maintainable, and extensible code. I wouldn't generally write a panic in production code. But, with coding challenges like this, using something that loudly interrupts a script to identify a false assumption is helpful
- There were a few cases where writing tests were helpful, but overall the quick and dirty method got me farther quicker 

## What needs improvement

### Stepping away from the computer
- This coding challenge was largely an exercise on the basics. One being that I need to step away from the computer more often when I'm working on something. I would have liked to have spent more time thinking about the problem, working with a whiteboard, or just generally resetting to come back fresher
> "All work and no play makes Jack a dull boy" - some smart person, probably not Jack


### Simplify and address base cases
- Sometimes I zoomed through without clearly addressing or thinking about the base cases... ¯\\_(ツ)_/¯
- Sitting down and actually writing out or visualizing cases were imperative for me for certain challenges

### Not skipping ahead sooner
- I started to try and catch up by doing parts 1 and 2 for every day even when they passed. I think I got more out of it when I switched over to just trying to catch up with part one for each day

## Things to try
### Test structure over a program structure
- Setting up challenges in the future to run via a test framework instead of a program.
  - The reason to do this to me would be to avoid breaking shared code used in part one of a problem when writing part two. I think this either almost, or did happen once
  - For things like this where I'm the only contributor and consumer of the code
    - I don't really need a main program
    - I don't want to invest too heavily in a test suite in addition to a main program, as I'm trying to move fast
    - So, committing to solving these problems via a test suite instead of a main program seems to make sense
  - Automating test runs via git hooks is straightforward, so this would be an easy win in regard to avoiding breakages and solving challenges

### Pomodoro Technique
  - I feel out my breaks, and take them at "sensible" times given the work and my focus. This can spiral into solving a problem that was outside of my initial goal. I think the [Pomodoro Technique](https://en.wikipedia.org/wiki/Pomodoro_Technique) would be an interesting constraint to try out

## Other Learnings
- I like that you have to set up a project and address whatever framework / language startup challenges there might be. If you are in a larger codebase for a long time and haven't set up a project from scratch in a while, your bootstrapping skills might be rusty

---
# Results
## Legend
| Status              | Emoji     |
|---------------------|-----------|
| Completed           | ✅        |
| Tried but failed    | ❌        |
| Completed but after December  | ⏰        |
| Not attempted       | ❓        |


|           | Day 1    | Day 2    | Day 3    | Day 4    | Day 5    | Day 6    | Day 7    | Day 8    | Day 9    | Day 10   | Day 11   | Day 12   |
|-----------|----------|----------|----------|----------|----------|----------|----------|----------|----------|-----------|-----------|-----------|
| part one    | ✅       | ✅       | ✅       | ✅       | ✅       | ✅       | ✅       | ✅       | ✅       | ✅        | ✅        | ⏰        |
| part two    | ✅     | ✅  | ✅     | ✅     | ✅     |    ❓     | ❓  | ❌     | ❓     | ❓      | ❓         | ❌

# Specific learnings

## 021223_cubes - Premature Optimization
- Here I thought that setting up an approach in part one that could be fairly flexible would set me up for success in part two. This was generally kind of true, and I did have fun using polymorphism and generics. But, it was kind of hit or miss to think this way for most of the other challenges, and made me think of the below quote: 
> Premature optimization is the root of all evil (or at least most of it) in programming. - [Donald Knuth](https://en.wikipedia.org/wiki/Donald_Knuth)

## 051223_seed_map - Different approaches
- This was fun, because I tried solving it a few ways. I solved part one the brute force way. When I realized this wouldn't cut it for part two I figured it would be nice to use one of Go's strengths, concurrency, given this is all CPU bound work. But, I had the sliding window / bounds approach in my head as well so I fleshed that idea out too.
- I wanted to see what other folks had done after I did those two approaches. When I went to the subreddit I found an approach where you work backwards from the given locations and see if it's a valid seed. I wouldn't have thought of that at the time, so it was a nice reminder to approach some problems like that

## 121223_cards - Taking a break
- I remember starting part one in the last days of December, getting kind of stuck, then feeling like I'd run out of time and stopping. When I picked this up to clean it up I couldn't resist trying again, and having fresh eyes was helpful
- I ended up simplifying my approach part one to using a brute force method, and I couldn't resist starting part two. It took writing down the base cases in text files for me to really understand how to calculate results without the brute force method. I wished I had done that sooner
- For part two I started to do something similar to my initial approach in grouping the possible combinations. I think I was on the right track. Part of my approach was using a cache and memoization to calculate results, and it seems other folks did that as well to get over the scalability issue. In the end I stopped before I could reliably group the possible combinations, but c'est la vie