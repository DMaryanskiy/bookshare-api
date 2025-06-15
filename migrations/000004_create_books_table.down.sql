DROP TABLE IF EXISTS books.books;

DROP TRIGGER IF EXISTS set_updated_at_books_trigger ON books.books;

DROP FUNCTION IF EXISTS books.set_updated_at();

DROP SCHEMA IF EXISTS books;
