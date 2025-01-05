CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	email VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE todos (
	id SERIAL PRIMARY KEY,
	title VARCHAR(100) NOT NULL,
	user_id INT NOT NULL,
	is_completed BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id),
)

INSERT INTO users (name, email) VALUES ('Naruto Uzumaki', 'naruto@gmail.com');

INSERT INTO todos (title, user_id) VALUES ('Buy groceries', 5);
INSERT INTO todos (title, user_id) VALUES ('Clean kunai', 5);
INSERT INTO todos (title, user_id) VALUES ('Practice shadow clone', 5);
INSERT INTO todos (title, user_id) VALUES ('Buy flowers', 5);
INSERT INTO todos (title, user_id) VALUES ('Chakra training', 5);

ALTER TABLE todos ADD COLUMN completed_at TIMESTAMP DEFAULT NULL;

UPDATE todos SET is_completed = TRUE WHERE id = 1;

SELECT * FROM todos WHERE user_id = 2 ORDER BY created_at;

ALTER TABLE todos ADD COLUMN due_date TIMESTAMP DEFAULT (CURRENT_TIMESTAMP + INTERVAL '7 days');

SELECT  * FROM todos WHERE due_date < CURRENT_TIMESTAMP AND is_completed = false ORDER BY due_date;

SELECT user_id, count(id) FROM todos GROUP BY user_id;

ALTER TABLE todos ADD column description VARCHAR(500);

SELECT constraint_name
FROM information_schema.table_constraints
WHERE table_name = 'todos' AND constraint_type = 'FOREIGN KEY';

ALTER TABLE todos DROP CONSTRAINT todos_user_id_fkey;

ALTER TABLE todos ADD CONSTRAINT todos_user_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

DELETE FROM users WHERE id = 1;

SELECT u.name, t.title, t.created_at 
FROM users u 
JOIN todos t ON u.id = t.user_id 
WHERE t.created_at = (
	SELECT max(created_at) 
	FROM todos 
	WHERE todos.user_id = t.user_id
);

SELECT name, title, created_at
FROM (
    SELECT u.name, t.title AS latest_todo, t.created_at,
           ROW_NUMBER() OVER (PARTITION BY t.user_id ORDER BY t.created_at DESC) AS rn
    FROM users u
    JOIN todos t ON u.id = t.user_id
) subquery
WHERE rn = 1;

SELECT name, subquery1.email, completed_count, pending_count
FROM (
    SELECT u.name, u.email,
           COUNT(t.id) AS completed_count
    FROM users u
    JOIN todos t ON u.id = t.user_id
    WHERE t.is_completed = true
    GROUP BY u.name, u.email
) subquery1
JOIN
(
	SELECT u.email,
           COUNT(t.id) AS pending_count
    FROM users u
    JOIN todos t ON u.id = t.user_id
    WHERE t.is_completed = false
    GROUP BY u.name, u.email
) subquery2
on subquery1.email = subquery2.email

create type todo_status_enum as ENUM('Completed','Pending','In Progress');

alter table todos add column status todo_status_enum not null default 'Pending';

update todos set status = case
	when is_completed = true then 'Completed'::todo_status_enum
	when is_completed = false then 'Pending'::todo_status_enum
	else 'Pending'::todo_status_enum
end;

alter table todos drop column is_completed;