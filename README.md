# MUR_Go
Go port of my Marvel Unlimited comic downloader, Marvel Universe Ripper.
[Windows binaries](https://github.com/Sorrow446/MUR_Go/releases)

![](https://i.imgur.com/xbfNr6j.png)

# Setup
**A subscription is required.**  
1. Login to [Marvel](https://dereferer.me/?http%3A//www.marvel.com/).
2. Install [EditThisCookie Chrome extension](https://chrome.google.com/webstore/detail/editthiscookie/fngmhnnpilhplaeedifhccceomclgfbg?hl=en) (any other Netscape extensions will also work).
3. Dump cookies to txt file named "cookies.txt" (http://www.marvel.com/ tab only).
4. Move cookies to MUR's directory.

Cookies will eventually expire. Just repeat the dumping process.

# Usage
All comics pages from Marvel Unlimited come in jpg format. They'll be converted losslessly to CBZ.

Download two comics:   
`mur_x86.exe https://www.marvel.com/comics/issue/89543/captain_america_2018_28 https://read.marvel.com/#/book/56459`

Download a single comic and from two text files:   
`mur_x64.exe https://www.marvel.com/comics/issue/89543/captain_america_2018_28 G:\1.txt G:\2.txt`

If building or running from source, you'll need to include the structs.   
`go run main.go structs.go <urls>...`
