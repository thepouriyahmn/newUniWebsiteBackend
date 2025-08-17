-- اضافه کردن ستون term_id به جدول classes
ALTER TABLE classes
ADD COLUMN term_id INT,
ADD CONSTRAINT fk_classes_term
FOREIGN KEY (term_id) REFERENCES terms(id)
ON DELETE CASCADE;

-- اضافه کردن ستون term_id به جدول users_classes
ALTER TABLE users_classes
ADD COLUMN term_id INT,
ADD CONSTRAINT fk_users_classes_term
FOREIGN KEY (term_id) REFERENCES terms(id)
ON DELETE CASCADE;