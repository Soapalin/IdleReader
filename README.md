# IdleReader: A Text-based Journey 📚📖
Welcome to Idle Reader! Immerse yourself in this idle text-based game! Accumulate knowledge, grow your intelligence to purchase real in-game books. 
<br>
To play the game, head to the  <a href="#-installation">installation</a> section down below or to the <a href="https://github.com/Soapalin/IdleReader/releases">release</a> page.

## 📕 Table of Contents


1. <a href="#-description">Description </a>
1. <a href="#-features">Features </a>
1. <a href="#-installation">Installation </a>
1. <a href="#-wiki">Wiki </a>
1. <a href="#-roadmap">Roadmap </a>
1. <a href="#-license">License </a>


## 📕 Description 

Start your journey with only your favourite book in-hand! By reading it, grow your Intelligence and Knowledge. Knowledge is the currency used in-game: use it wisely to purchase and collect your favourite books! Beware, some books have Intelligence requirements, and the only way to grow your IQ is to read different books (at least once). 

Learn more about the <a href="#features">features</a> below or in explore our <a href="https://github.com/Soapalin/IdleReader/wiki">wiki</a> page for in-depth game guides.


<img src="IdleReader_demo.gif" />


## 📕 Features
- 📗 Idle Reading: Leave the game in the background and continue progressing
- 📗 Real Books: Purchase real in-game books with the knowledge you have accumulated
- 📗 Bookshop: Refreshed regularly for you to find your favourite books
- 📗 IQ Growth: Reading different books will increase your IQ!
- 📗 Knowledge Growth: Select different books for different rate of growth
- 📗 Concurrent Reading: Grow your IQ high enough, and read up to 3 books simultaneously!
- 📗 Auction: Tired of waiting for your favourite books to be available in the bookshop? Head to the Auction and bid for it! Beware, prices are not low...


## 📕 Installation
For a stress-free installation, head to release page and download the latest version of the game. Only the executable `IdleReader` file is required to run the game. A version is available for each OS (Windows, Linux, MacOS).

To play the game, simply double click on the file that you have just downloaded.

Note for Windows users: For best experience, install `Windows Terminal` (default as of Windows 11 22H2!). See in-depth guide <a href="https://github.com/Soapalin/IdleReader/wiki/Installing-Windows-Terminal">here</a>. You can still play the game without `Windows Terminal` but the default rendering will be a bit off (such as not having emojis available).


For advanced users, you can build your own binaries of the game:
- Clone the repository: `git clone [repository]`
- Navigate to the repo and build the binaries with `go build .`
- the binaries built should be called `engine`


## 📕 Wiki

The game wiki and FAQs can be found here:

<a href="https://github.com/Soapalin/IdleReader/wiki">IdleReader Wiki</a>

This includes in-depth game guides and tutorials.

## 📕 Roadmap

Current version: v0.0.1

### Road to v0.1

- [x] AllBooksLibrary migrating to SQLite
- [x] Exit View
- [x] Create adequate README.md
- [x] Create wiki
- [x] Update wiki on unicode rendering (use Windows Terminal if still using cmd.exe)
- [x] Store saves/database in different path depending on OS
- [x] Create a release page for pre-built .exe

### Road to v0.2
- [x] Auction feature with higher price
- [x] Search Book function in Auction
- [x] Shop shows books depending on your IQ
- [x] Update Help View
- [x] Add a way to display whether item has already been purchased or not in Bookshop

### Road to v0.3
- [x] Add item effects
- [x] Search Book function in My Bookshelf
- [x] Add cost to auction search
- [x] Add price algorithm
- [x] Fetch and add real books (use external API or web scraper)
- [x] Algorithm with reading speed

### Road to v1.0 
- [] Find solution backwards compatibility with PlayerSave gob 
- [] User Testing needed

### Road to v2.0 
- [] Find gameplay for prestige 



## 📕 License
This project is licensed under the terms of the MIT License. See the LICENSE file for details.

## 📕 Closing Thoughts & How To Contribute
Thank you for checking out my text-based idle game, IdleReader!

Contributions are welcome! Since it is my first full game that I have decided to release, I am still pretty new to this and there are many ways to contribute to Idle Reader. Here are a few asks:
- Book Recommendations: Suggest books to add to the game's library, covering a wide range of genres and topics.
- Pull Requests: Help improve the game's functionality, optimize performance, and fix any bugs or issues.
- New Ideas: Share or code new features and ideas and don't hesitate to reach out to me directly!


Till Next Time
