# The Challenge - Building wc

The functional requirements for wc are concisely described by it’s man page - give it a go in your local terminal now:
man wc
The TL/DR version is: wc – word, line, character, and byte count.
Step Zero
Like all good software engineering we’re zero indexed! In this step you’re going to set your environment up ready to begin developing and testing your solution.
I’ll leave you to setup your IDE / editor of choice and programming language of choice. After that here’s what I’d like you to do to be ready to test your solution.
Download the following text: https://www.gutenberg.org/cache/epub/132/pg132.txt and save it as text.txt

## Step One
In this step your goal is to write a simple version of wc, let’s call it ccwc (cc for Coding Challenges) that takes the command line option -c and outputs the number of bytes in a file.
If you’ve done it right your output should match this:
>ccwc -c test.txt
  341836 test.txt
If it doesn’t, check your code, fix any bugs and try again. If it does, congratulations! On to…

## Step Two
In this step your goal is to support the command line option -l that outputs the number of lines in a file.
If you’ve done it right your output should match this:
>ccwc -l test.txt
    7137 test.txt
If it doesn’t, check your code, fix any bugs and try again. If it does, congratulations! On to…

## Step Three
In this step your goal is to support the command line option -w that outputs the number of words in a file. If you’ve done it right your output should match this:
>ccwc -w test.txt
   58159 test.txt
If it doesn’t, check your code, fix any bugs and try again. If it does, congratulations! On to…

## Step Four
In this step your goal is to support the command line option -m that outputs the number of characters in a file. If the current locale does not support multibyte characters this will match the -c option.
You can learn more about programming for locales here: https://learn.microsoft.com/en-us/globalization/locale/locale-and-culture
For this one your answer will depend on your locale, so if can, use wc itself and compare the output to your solution:
>wc -m test.txt
  339120 test.txt

>ccwc -m test.txt
  339120 test.txt
If it doesn’t, check your code, fix any bugs and try again. If it does, congratulations! On to…
Step Five
In this step your goal is to support the default option - i.e. no options are provided, which is the equivalent to the -c, -l and -w options. If you’ve done it right your output should match this:
>ccwc test.txt          
    7137   58159  341836 test.txt
If it doesn’t, check your code, fix any bugs and try again. If it does, congratulations! On to…
The Final Step
In this step your goal is to support being able to read from standard input if no filename is specified. If you’ve done it right your output should match this:
>cat test.txt | ccwc -l 
    7137
If it doesn’t, check your code, fix any bugs and try again. If it does, congratulations! You’ve done it, pat yourself on the back, job well done!