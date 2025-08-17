ALTER TABLE classes
DROP FOREIGN KEY fk_classes_term,
DROP COLUMN term_id;

-- حذف foreign key و ستون از جدول users_classes
ALTER TABLE users_classes
DROP FOREIGN KEY fk_users_classes_term,
DROP COLUMN term_id;