DELETE FROM `auth` WHERE id IN (1, 2);

--bun:split

DELETE FROM `user` WHERE id IN (1, 2, 3);
