ALTER TABLE comments
  ADD CONSTRAINT comments_post_id_fkey
    FOREIGN KEY (post_id) REFERENCES  posts(id) ON DELETE CASCADE,
  ADD CONSTRAINT comments_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

