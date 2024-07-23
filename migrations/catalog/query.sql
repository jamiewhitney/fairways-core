-- name: GetACourse :one
select *
from courses
where id = ?
limit 1;

-- name: ListCourses :many
SELECT *
FROM courses
ORDER BY id
limit ? offset ?;

-- name: CreateCourse :execlastid
insert into courses (name, holes, lapsed, live, street_address_1, city, state, postal_code, country, latitude,
                     longitude)
values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);