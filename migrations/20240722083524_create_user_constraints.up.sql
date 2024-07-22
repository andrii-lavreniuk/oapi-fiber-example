ALTER TABLE user_data
  ADD CONSTRAINT user_data_user_FK
  FOREIGN KEY (user_id)
  REFERENCES user(id)
  ON DELETE CASCADE
  ON UPDATE NO ACTION;

--bun:split

ALTER TABLE user_profile
  ADD CONSTRAINT user_profile_user_FK
  FOREIGN KEY (user_id)
  REFERENCES user(id)
  ON DELETE CASCADE
  ON UPDATE NO ACTION;


