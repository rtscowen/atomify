# atomify
A minimalist habit tracker, based on the principles of "Atomic Habits" by James Clear. For every habit you create, 
define how many times a day you want to complete it. Get a daily todo list. Take break days. 

## Installation

From releases, download the appropriate binary for your system, and rename to atomify

## Usage

```
atomify new gym 1
```

Creates a new habit called "gym" with an expectation of once per day

```
atomify remove gym 
```

Deletes a habit by name

```
atomify list
```

List all the habits you are currently tracking 

```
atomify day
```

Get your todo list for the day, comprised of habits you haven't hit your daily cap for yet

```
atomify ++ gym
```

Do this when you've been to the gym. Marks one unit of the "gym" habit, and so in this case removes it from the todo list, as it's only once per day.

```
atomoify -- gym
```

Undoes the above command.

``` 
atomify break gym
```

Have a break from this habit. Removes it from the todo list without incrementing the count. 

```
atomoify unbreak gym
```

Undoes the break command

```
atomify update gym -r=[new name] -d=[new daily expectation]
```

Update is used to rename a habit and/or change the daily expectation. Flags can be used independently, as long as one is used. 
