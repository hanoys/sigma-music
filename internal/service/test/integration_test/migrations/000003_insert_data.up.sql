INSERT INTO users(id, name, email, phone, password, salt, country)
VALUES ('1add32df-d439-4fd1-9d4c-bef946b4a1fc', 'Timur', 'timur@mail.ru', '+79999999999', 'passwd', 'salt', 'Russia');

INSERT INTO musicians(id, name, email, password, salt, country, description)
VALUES ('1add32df-d439-4fd1-9d4c-bef946b4a1fa', 'Timur', 'timur@mail.ru', 'passwd', 'salt', 'Russia', 'description');

INSERT INTO albums(id, name, description, published, release_date)
VALUES ('b24fa8eb-9df6-406c-9b45-763d7b5a5078', 'album', 'description', true, now());

INSERT INTO album_musician(musician_id, album_id) 
VALUES ('1add32df-d439-4fd1-9d4c-bef946b4a1fa', 'b24fa8eb-9df6-406c-9b45-763d7b5a5078');

INSERT INTO tracks(id, album_id, name, url)
VALUES ('41623ac1-b98d-4478-a10f-870a80c697b6', 'b24fa8eb-9df6-406c-9b45-763d7b5a5078', 'trackname', 'url');

INSERT INTO genres(id, name)
VALUES ('32f24dfc-3823-41e4-a073-c3553c981db1', 'genre');

INSERT INTO track_genre(track_id, genre_id)
VALUES('41623ac1-b98d-4478-a10f-870a80c697b6', '0dba1d8d-cbb4-4126-a5b2-6596866a2d7a');

INSERT INTO favorite(user_id, track_id)
VALUES ('1add32df-d439-4fd1-9d4c-bef946b4a1fc', '41623ac1-b98d-4478-a10f-870a80c697b6');

INSERT INTO comments(id, user_id, track_id, stars, comment_text)
VALUES ('cad092f7-d0d3-4a9a-83ae-b0d12f4bb382', '1add32df-d439-4fd1-9d4c-bef946b4a1fc', '41623ac1-b98d-4478-a10f-870a80c697b6', 3, 'text');

INSERT INTO users_history(id, user_id, track_id) 
VALUES ('cad092f7-d0d3-4a9a-83ae-b0d12f4bb383', '1add32df-d439-4fd1-9d4c-bef946b4a1fc', '41623ac1-b98d-4478-a10f-870a80c697b6');

