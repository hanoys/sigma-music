-- genres
INSERT INTO genres(id, name)
VALUES ('c783ef09-4444-402d-8f9a-9b766dc8e2f9', 'rock');

INSERT INTO genres(id, name)
VALUES ('f3ef318d-dc9a-4c1d-999f-7aff3fca250e', 'rap');

INSERT INTO genres(id, name)
VALUES ('518bc817-bb59-4616-8310-e19b415d2a5d', 'pop');

-- users
INSERT INTO users(id, name, email, phone, password, country)
VALUES ('8e213888-93e9-4a27-af95-aebaab37103c', 'Timur', 'timur@mail.ru', '1', 'password1', 'russia');

INSERT INTO users(id, name, email, phone, password, country)
VALUES ('7877569c-a853-475c-80a9-8824865fe6de', 'Pasha', 'pasha@mail.ru', '2', 'password2', 'russia');

INSERT INTO users(id, name, email, phone, password, country)
VALUES ('64422934-7306-4045-bce4-a64c9de80d7a', 'Emir', 'emir@mail.ru', '3', 'password3', 'russia');

-- musicians
INSERT INTO musicians(id, name, email, password, country, description)
VALUES ('bb56e25a-41b1-44fa-8c5b-bd62ab31e87d', 'Eminem', 'eminem@yandex.ru', 'passwordeminem', 'russia', 'eminen description');

INSERT INTO musicians(id, name, email, password, country, description)
VALUES ('ab934887-55c5-4bd3-9115-686a0681a901', 'musician', 'musician@mail.ru', 'pass', 'russia', 'musician description');

INSERT INTO musicians(id, name, email, password, country, description)
VALUES ('cbd2d7f3-fcda-41bd-adfa-54c3e1d61f57', 'muse', 'muse@gmail.com', 'word', 'russia', 'muse description');

-- albums
INSERT INTO albums(id, name, description, published, release_date)
VALUES ('a427eafd-5f5b-45a6-b4d5-ed32fe7df3d3', 'album1', 'descr1', true, now());

INSERT INTO albums(id, name, description, published, release_date)
VALUES ('c876fe34-f0dd-4a9d-8391-1c388050774a', 'album2', 'descr2', true, now());

INSERT INTO albums(id, name, description, published, release_date)
VALUES ('58159512-4721-4c58-9693-6a08b51fbd37', 'album3', 'descr3', true, now());

INSERT INTO albums(id, name, description, published, release_date)
VALUES ('b255f791-014d-4621-9555-ab9f182774cc', 'album4', 'descr4', true, now());

INSERT INTO albums(id, name, description, published, release_date)
VALUES ('6951baa3-51f2-4f3d-8a21-8b8cb25c5205', 'album5', 'descr5', true, now());

-- album_musician

INSERT INTO album_musician(musician_id, album_id)
VALUES ('bb56e25a-41b1-44fa-8c5b-bd62ab31e87d', 'a427eafd-5f5b-45a6-b4d5-ed32fe7df3d3');

INSERT INTO album_musician(musician_id, album_id)
VALUES ('bb56e25a-41b1-44fa-8c5b-bd62ab31e87d', 'c876fe34-f0dd-4a9d-8391-1c388050774a');

INSERT INTO album_musician(musician_id, album_id)
VALUES ('ab934887-55c5-4bd3-9115-686a0681a901', '58159512-4721-4c58-9693-6a08b51fbd37');

INSERT INTO album_musician(musician_id, album_id)
VALUES ('ab934887-55c5-4bd3-9115-686a0681a901', 'b255f791-014d-4621-9555-ab9f182774cc');

INSERT INTO album_musician(musician_id, album_id)
VALUES ('cbd2d7f3-fcda-41bd-adfa-54c3e1d61f57', '6951baa3-51f2-4f3d-8a21-8b8cb25c5205');

-- tracks and genres
INSERT INTO tracks(id, album_id, name, url)
VALUES ('548438c8-bb45-43eb-8588-2c748d6053d9', 'a427eafd-5f5b-45a6-b4d5-ed32fe7df3d3', 'track1', 'testurl');
INSERT INTO track_genre(track_id, genre_id)
VALUES ('548438c8-bb45-43eb-8588-2c748d6053d9', 'f3ef318d-dc9a-4c1d-999f-7aff3fca250e');

INSERT INTO tracks(id, album_id, name, url)
VALUES ('73362edc-0f16-4295-9e67-3796e2d36445', 'a427eafd-5f5b-45a6-b4d5-ed32fe7df3d3', 'track2', 'testurl');
INSERT INTO track_genre(track_id, genre_id)
VALUES ('73362edc-0f16-4295-9e67-3796e2d36445', 'f3ef318d-dc9a-4c1d-999f-7aff3fca250e');

INSERT INTO tracks(id, album_id, name, url)
VALUES ('4dcd45ab-c7e0-4362-8531-1cee1a23a3b3', 'c876fe34-f0dd-4a9d-8391-1c388050774a', 'track3', 'testurl');
INSERT INTO track_genre(track_id, genre_id)
VALUES ('4dcd45ab-c7e0-4362-8531-1cee1a23a3b3', 'f3ef318d-dc9a-4c1d-999f-7aff3fca250e');

INSERT INTO tracks(id, album_id, name, url)
VALUES ('fdd4c1d1-085d-46c4-af4f-3b931f6ab593', 'c876fe34-f0dd-4a9d-8391-1c388050774a', 'track4', 'testurl');
INSERT INTO track_genre(track_id, genre_id)
VALUES ('fdd4c1d1-085d-46c4-af4f-3b931f6ab593', 'f3ef318d-dc9a-4c1d-999f-7aff3fca250e');


--INSERT INTO tracks(id, album_id, name, url)
--VALUES ('6eb41931-fe77-408b-9aad-821524d37aae3', '58159512-4721-4c58-9693-6a08b51fbd37', 'track5', 'testurl');
--INSERT INTO track_genre(track_id, genre_id)
--VALUES ('6eb41931-fe77-408b-9aad-821524d37aae3', '518bc817-bb59-4616-8310-e19b415d2a5d');

INSERT INTO tracks(id, album_id, name, url)
VALUES ('71a65dc5-db81-402c-96ea-c8d149d710a4', '58159512-4721-4c58-9693-6a08b51fbd37', 'track6', 'testurl');
INSERT INTO track_genre(track_id, genre_id)
VALUES ('71a65dc5-db81-402c-96ea-c8d149d710a4', '518bc817-bb59-4616-8310-e19b415d2a5d');


INSERT INTO tracks(id, album_id, name, url)
VALUES ('762380ca-992b-47ec-af6f-4d9efc975ba5', 'b255f791-014d-4621-9555-ab9f182774cc', 'track7', 'testurl');
INSERT INTO track_genre(track_id, genre_id)
VALUES ('762380ca-992b-47ec-af6f-4d9efc975ba5', 'c783ef09-4444-402d-8f9a-9b766dc8e2f9');

INSERT INTO tracks(id, album_id, name, url)
VALUES ('b7b5e9e2-a80f-425d-a0af-6bc568bc9869', 'b255f791-014d-4621-9555-ab9f182774cc', 'track8', 'testurl');
INSERT INTO track_genre(track_id, genre_id)
VALUES ('b7b5e9e2-a80f-425d-a0af-6bc568bc9869', 'c783ef09-4444-402d-8f9a-9b766dc8e2f9');


INSERT INTO tracks(id, album_id, name, url)
VALUES ('6e3f562e-5907-42a7-80fd-2edfcbbeda2a', '6951baa3-51f2-4f3d-8a21-8b8cb25c5205', 'track9', 'testurl');
INSERT INTO track_genre(track_id, genre_id)
VALUES ('6e3f562e-5907-42a7-80fd-2edfcbbeda2a', 'c783ef09-4444-402d-8f9a-9b766dc8e2f9');

INSERT INTO tracks(id, album_id, name, url)
VALUES ('8ec5a780-1f8d-41d3-b032-bb3c668a8825', '6951baa3-51f2-4f3d-8a21-8b8cb25c5205', 'track10', 'testurl');
INSERT INTO track_genre(track_id, genre_id)
VALUES ('8ec5a780-1f8d-41d3-b032-bb3c668a8825', 'c783ef09-4444-402d-8f9a-9b766dc8e2f9');

-- users_history

INSERT INTO users_history(user_id, track_id)
VALUES ('8e213888-93e9-4a27-af95-aebaab37103c', '548438c8-bb45-43eb-8588-2c748d6053d9');

INSERT INTO users_history(user_id, track_id)
VALUES ('8e213888-93e9-4a27-af95-aebaab37103c', '73362edc-0f16-4295-9e67-3796e2d36445');

INSERT INTO users_history(user_id, track_id)
VALUES ('8e213888-93e9-4a27-af95-aebaab37103c', '4dcd45ab-c7e0-4362-8531-1cee1a23a3b3');

INSERT INTO users_history(user_id, track_id)
VALUES ('8e213888-93e9-4a27-af95-aebaab37103c', 'fdd4c1d1-085d-46c4-af4f-3b931f6ab593');

INSERT INTO users_history(user_id, track_id)
VALUES ('8e213888-93e9-4a27-af95-aebaab37103c', 'b7b5e9e2-a80f-425d-a0af-6bc568bc9869');

INSERT INTO users_history(user_id, track_id)
VALUES ('8e213888-93e9-4a27-af95-aebaab37103c', '548438c8-bb45-43eb-8588-2c748d6053d9');

INSERT INTO users_history(user_id, track_id)
VALUES ('7877569c-a853-475c-80a9-8824865fe6de', '548438c8-bb45-43eb-8588-2c748d6053d9');

