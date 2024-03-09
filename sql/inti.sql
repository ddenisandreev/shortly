create schema shortly;

create table shortly.urls (
    id      serial  PRIMARY KEY, 
    url_    text    NOT NULL
);

insert into shortly.urls (url_) values ('http://google.com')

SELECT * FROM shortly.urls