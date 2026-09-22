ALTER TABLE comments
  DROP CONSTRAINT comments_post_id_fkey,
  DROP CONSTRAINT comments_user_id_fkey;
