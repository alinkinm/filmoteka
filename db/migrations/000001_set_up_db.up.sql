CREATE TABLE Actors (
    id serial primary key,
    names varchar not null,
    sex char(1) not null, 
    bd Date not null,
    UNIQUE(names)
);

CREATE TABLE Movies (
    id serial primary key,
    title varchar(150) not null,
    descr varchar(1000) not null, 
    release Date not null,
    rating int not null,
    CHECK (title <> ''),
    UNIQUE(title)
);

CREATE TABLE ActorMovie (
    actor_id int references Actors (id),
    movie_id int references Movies (id)
);

CREATE OR REPLACE FUNCTION add_actors_to_movie(actor_ids integer[], movie_id int) RETURNS VOID AS $$
BEGIN
    INSERT INTO ActorMovie (actor_id, movie_id)
    SELECT actor_id, movie_id
    FROM unnest(actor_ids) AS actor_id
    WHERE EXISTS (
        SELECT 1
        FROM Actors a
        WHERE a.id = actor_id
    );
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_actors_from_movie(actor_ids integer[], movie_id int) RETURNS VOID AS $$
BEGIN
    DELETE FROM ActorMovie
    WHERE actor_id = ANY(actor_ids)
    AND movie_id = movie_id;
END;
$$ LANGUAGE plpgsql;

CREATE INDEX idx_movies_title ON Movies(title);

CREATE INDEX idx_actors_name ON Actors(names);

