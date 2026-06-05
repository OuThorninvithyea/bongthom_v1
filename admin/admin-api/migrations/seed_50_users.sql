-- Seed 50 mock users into tbl_users
-- Password for all: "password123" (bcrypt hashed)
-- Run: psql "postgresql://outhorninvuth:postgres@127.0.0.1:5432/postgres?sslmode=disable" -f migrations/seed_50_users.sql

BEGIN;

INSERT INTO tbl_users (first_name, last_name, user_name, email, password, role_name, role_id, is_admin, login_session, last_login, currency_id, language_id, status_id, "order", created_by, created_at)
VALUES
('Sokha',     'Chea',       'MOCK_USER_01', 'mock01@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Admin',      1, true,  NULL, NULL, 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Sopheap',   'Meas',       'MOCK_USER_02', 'mock02@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Sreyneang', 'Sok',        'MOCK_USER_03', 'mock03@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 1, 2, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Dara',      'Heng',       'MOCK_USER_04', 'mock04@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 2, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Bopha',     'Ly',         'MOCK_USER_05', 'mock05@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 1, 3, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Rithy',     'Oum',        'MOCK_USER_06', 'mock06@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Admin',      1, true,  NULL, NULL, 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Sophea',    'Chhum',      'MOCK_USER_07', 'mock07@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 2, 2, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Vannak',    'Kong',       'MOCK_USER_08', 'mock08@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 1, 1, 2, nextval('tbl_users_id_seq'), 1, NOW()),
('Sokunthea', 'Touch',      'MOCK_USER_09', 'mock09@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Dalin',     'Phan',       'MOCK_USER_10', 'mock10@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 1, 1, 3, nextval('tbl_users_id_seq'), 1, NOW()),
('Kosal',     'Seng',       'MOCK_USER_11', 'mock11@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 2, 2, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Sreymom',   'Khim',       'MOCK_USER_12', 'mock12@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 3, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Chamroeun', 'Yin',        'MOCK_USER_13', 'mock13@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Sovann',    'Huot',       'MOCK_USER_14', 'mock14@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 1, 2, 2, nextval('tbl_users_id_seq'), 1, NOW()),
('Sokchea',   'Pov',        'MOCK_USER_15', 'mock15@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 1, 3, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Maly',      'Sam',        'MOCK_USER_16', 'mock16@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 2, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Piseth',    'Um',         'MOCK_USER_17', 'mock17@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NOW(), 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Sovath',    'Mam',        'MOCK_USER_18', 'mock18@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Chanthou',  'Sorn',       'MOCK_USER_19', 'mock19@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 3, 2, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Bunthoeun','Prak',       'MOCK_USER_20', 'mock20@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Rothana',   'Suon',       'MOCK_USER_21', 'mock21@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Admin',      1, true,  NULL, NOW(), 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Sokhom',    'Hul',        'MOCK_USER_22', 'mock22@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 2, 3, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Sothea',    'Say',        'MOCK_USER_23', 'mock23@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Sophat',    'Keo',        'MOCK_USER_24', 'mock24@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 1, 2, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Rany',      'Tep',        'MOCK_USER_25', 'mock25@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 1, 1, 2, nextval('tbl_users_id_seq'), 1, NOW()),
('Sokheng',   'Nop',        'MOCK_USER_26', 'mock26@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Sokmean',   'In',         'MOCK_USER_27', 'mock27@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NOW(), 2, 2, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Sreyroth',  'Chhorn',     'MOCK_USER_28', 'mock28@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Dany',      'Chuon',      'MOCK_USER_29', 'mock29@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 3, 3, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('John',      'Smith',      'MOCK_USER_30', 'mock30@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 1, 2, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Mary',      'Johnson',    'MOCK_USER_31', 'mock31@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('William',   'Brown',      'MOCK_USER_32', 'mock32@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 2, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Linda',     'Davis',      'MOCK_USER_33', 'mock33@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Admin',      1, true,  NULL, NOW(), 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Michael',   'Miller',     'MOCK_USER_34', 'mock34@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 1, 3, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Patricia',  'Wilson',     'MOCK_USER_35', 'mock35@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 3, 1, 2, nextval('tbl_users_id_seq'), 1, NOW()),
('Robert',    'Moore',      'MOCK_USER_36', 'mock36@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 1, 2, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Jennifer',  'Taylor',     'MOCK_USER_37', 'mock37@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 2, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('David',     'Anderson',   'MOCK_USER_38', 'mock38@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Sarah',     'Thomas',     'MOCK_USER_39', 'mock39@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 1, 3, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('James',     'Jackson',    'MOCK_USER_40', 'mock40@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Admin',      1, true,  NULL, NULL, 2, 2, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Barbara',   'White',      'MOCK_USER_41', 'mock41@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Richard',   'Harris',     'MOCK_USER_42', 'mock42@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NOW(), 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Susan',     'Martin',     'MOCK_USER_43', 'mock43@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 3, 1, 3, nextval('tbl_users_id_seq'), 1, NOW()),
('Joseph',    'Thompson',   'MOCK_USER_44', 'mock44@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 1, 2, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Jessica',   'Garcia',     'MOCK_USER_45', 'mock45@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 2, 1, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Thomas',    'Martinez',   'MOCK_USER_46', 'mock46@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 1, 3, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Lisa',      'Robinson',   'MOCK_USER_47', 'mock47@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 1, 1, 2, nextval('tbl_users_id_seq'), 1, NOW()),
('Charles',   'Clark',      'MOCK_USER_48', 'mock48@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Operator',   2, false, NULL, NULL, 1, 2, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Nancy',     'Rodriguez',  'MOCK_USER_49', 'mock49@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Admin',      1, true,  NULL, NOW(), 2, 2, 1, nextval('tbl_users_id_seq'), 1, NOW()),
('Daniel',    'Lewis',      'MOCK_USER_50', 'mock50@bongthom.com', '$2a$10$Ee8w8i1Uc9wi3gzioEgWm.rbXXhf3UOwRDpyP09MLnzHl0.Mk2CKO', 'Viewer',     3, false, NULL, NULL, 1, 1, 1, nextval('tbl_users_id_seq'), 1, NOW())
ON CONFLICT (user_name) DO NOTHING;

COMMIT;
