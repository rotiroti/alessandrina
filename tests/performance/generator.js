"use strict";

const TITLE = [
  "The Go Programming Language",
  "Patterns of Enterprise Application Architecture",
  "Clean Code",
  "Clean Architecture",
  "The Pragmatic Programmer",
  "Design Patterns",
  "Code Complete",
  "Refactoring",
  "Working Effectively with Legacy Code",
  "Domain-Driven Design",
  "The Mythical Man-Month",
  "Extreme Programming Explained",
  "The Art of Computer Programming",
  "Introduction to Algorithms",
  "Code",
  "Serverless Architectures on AWS",
  "Cloud Computing - Concepts, Technology & Architecture",
  "Docker in Action",
  "Microservices Patterns",
  "Designing Data-Intensive Applications",
];

const AUTHORS = [
  "Brian W. Kernighan, Alan Donovan",
  "Martin Fowler",
  "Robert C. Martin",
  "Andy Hunt, Dave Thomas",
  "Erich Gamma, Richard Helm, Ralph Johnson, John Vlissides",
  "Steve McConnell",
  "Martin Fowler",
  "Michael Feathers",
  "Eric Freeman, Bert Bates, Kathy Sierra, Elisabeth Robson",
  "Eric Evans",
  "Frederick P. Brooks Jr.",
  "Kent Beck, Cynthia Andres",
  "Donald E. Knuth",
  "Thomas H. Cormen, Charles E. Leiserson, Ronald L. Rivest, Clifford Stein",
  "Charles Petzold",
  "Peter Sbarski",
  "Thomas Erl",
  "Jeff Nickoloff",
  "Chris Richardson",
  "Martin Kleppmann",
];

const ISBN = [
  "9780134190440",
  "9780321127426",
  "9780132350884",
  "9780137081073",
  "9780201633610",
  "9780201485677",
  "9780735619678",
  "9780134757599",
  "9780132350884",
  "9780321127426",
  "9780134190440",
  "9780321127426",
  "9780134190440",
  "9780321127426",
  "9780134190440",
  "9780321127426",
  "9780134190440",
  "9780321127426",
  "9780134190440",
  "9780321127426",
];

const PUBLISHER = [
  "Addison-Wesley Professional",
  "Prentice Hall",
  "Microsoft Press",
  "O'Reilly Media",
  "Wiley",
  "Manning Publications",
  "No Starch Press",
  "Apress",
  "Pragmatic Bookshelf",
];

const generateRandomData = (context, events, done) => {
    const title = TITLE[Math.floor(Math.random() * TITLE.length)];
    const authors = AUTHORS[Math.floor(Math.random() * AUTHORS.length)];
    const isbn = ISBN[Math.floor(Math.random() * ISBN.length)];
    const publisher = PUBLISHER[Math.floor(Math.random() * PUBLISHER.length)];
    const pages = Math.floor(Math.random() * 1000);

    context.vars.title = title;
    context.vars.authors = authors;
    context.vars.isbn = isbn;
    context.vars.publisher = publisher;
    context.vars.pages = pages;

    return done()
}

module.exports = {
  generateRandomData
}
