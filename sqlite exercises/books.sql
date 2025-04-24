
-- .mode table or .mode box;
-- .read books.sql;
--
-- to enable foreign key support:
--   pragma foreign_keys = on;
-- use with:
--   foreign key (col) references table (col)


-- Clean up before running the script.

drop table if exists Books;
drop table if exists Authors;
drop table if exists Write;

-- Create tables.

create table Books (
  isbn text primary key,
  title text,
  year integer,
  pages integer
);

--create table Publish (
--  isbn text,
--  pubId integer
--);

create table Authors (
  authorId integer primary key,
  name text
);

create table Write (
  isbn string,
  authorId integer
);


-- Davies
insert into Authors values
   (1, 'Robertson Davies'),
   (2, 'Judith Skelton Grant');

insert into Books values
   ('0771595565', 'Tempest-Tost', 1951, 350),
   ('0771595566', 'Leaven of Malice', 1954, 340),
   ('0771595567', 'A Mixture of Frailties', 1958, 330),
   ('0771595568', 'Fifth Business', 1970, 410),
   ('0771595569', 'The Manticore', 1972, 420),
   ('0771595570', 'World of Wonders', 1975, 450),
   ('0771595571', 'The Rebel Angels', 1981, 390),
   ('0771595572', 'What''s Bred in the Bone', 1985, 490),
   ('0771595573', 'The Lyre of Orpheus', 1988, 420),
   ('0771595574', 'Murther and Walking Spirits', 1991, 290),
   ('0771595575', 'The Cunning Man', 1994, 310),
   ('0771333333', 'For Your Eye Alone: Letters 1976-1995', 1999, 610),
   ('0771333334', 'Discoveries: Early letters 1938-1975', 2002, 700);

insert into Write values
   ('0771595565', 1),
   ('0771595566', 1),
   ('0771595567', 1),
   ('0771595568', 1),
   ('0771595569', 1),
   ('0771595570', 1),
   ('0771595571', 1),
   ('0771595572', 1),
   ('0771595573', 1),
   ('0771595574', 1),
   ('0771595575', 1),
   ('0771333333', 1),
   ('0771333333', 2),
   ('0771333334', 1),
   ('0771333334', 2);

-- Woolfe
insert into Authors values
   (3, 'Virginia Woolf');

insert into Books values
   ('0771234567', 'Mrs Dalloway', 1925, 230),
   ('0771234568', 'To the Lighthouse', 1927, 210);

insert into Write values
   ('0771234567', 3),
   ('0771234568', 3);

-- Fowles
insert into Authors values
   (4, 'John Fowles');

insert into Books values
   ('0771123451', 'The Collector', 1963, 190),
   ('0771123452', 'The Magus', 1965, 720),
   ('0771123453', 'The French Lieutenant''s Woman', 1969, 360),
   ('0771123454', 'The Ebony Tower', 1974, 290),
   ('0771123455', 'Daniel Martin', 1977, 730),
   ('0771123456', 'Mantissa', 1982, 180),
   ('0771123457', 'A Maggot', 1985, 240);

insert into Write values
   ('0771123451', 4),
   ('0771123452', 4),
   ('0771123453', 4),
   ('0771123454', 4),
   ('0771123455', 4),
   ('0771123456', 4),
   ('0771123457', 4);

-- McMurtry
insert into Authors values
   (5, 'Larry McMurtry');

insert into Books values
   ('0771432151', 'The Last Picture Show', 1966, 210),
   ('0771432152', 'Texasville', 1987, 310),
   ('0771432153', 'Duane''s Depressed', 1999, 410),
   ('0771432154', 'When the Light Goes', 2007, 180),
   ('0771432155', 'Rhino Ranch', 2009, 200);

insert into Write values
   ('0771432151', 5),
   ('0771432152', 5),
   ('0771432153', 5),
   ('0771432154', 5),
   ('0771432155', 5);


------------------------------------------------------------
-- To test joins:

drop table if exists T;
drop table if exists U;
create table T (name text, j integer);
create table U (j integer, name text);
insert into T values ('A', 1), ('B', 2), ('C', 3), ('D', null);
insert into U values (1, 'W'), (1, 'X'), (3, 'Y'), (null, 'Z');
