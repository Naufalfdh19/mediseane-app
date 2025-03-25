create table logistics (
   id bigserial primary key,
   name varchar not null,
   price decimal(14, 2) not null,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);


CREATE TABLE postal_ongkir (
id bigserial PRIMARY KEY,
id_location int NOT NULL UNIQUE,
postal_code int NOT NULL UNIQUE,
created_at timestamp not null default current_timestamp,
updated_at timestamp not null default current_timestamp,
deleted_at timestamp null
);


create table partners (
   id bigserial primary key,
   name varchar not null,
   year_founded varchar not null,
   active_days varchar not null,
   start_hour TIME not null,
   end_hour TIME not null,
   is_active boolean not null,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);




create table product_categories (
   id bigserial primary key,
   name varchar not null,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);




create table product_classifications (
   id bigserial primary key,
   name varchar not null,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);




create table product_forms (
   id bigserial primary key,
   name varchar not null,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);


create table users (
   id bigserial primary key,
   name varchar not null,
   email varchar not null,
   password varchar not null,
   profile_photo varchar not null,
   role varchar not null,
   is_verified bool not null,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);


CREATE TABLE provinces (
   id BIGSERIAL PRIMARY KEY,
   name VARCHAR NOT NULL,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);


CREATE TABLE cities (
   id BIGSERIAL PRIMARY KEY,
   province_id BIGINT NOT NULL REFERENCES provinces(id),
   name VARCHAR NOT NULL   ,
   location GEOGRAPHY(Point, 4326) NULL,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);


CREATE TABLE districts (
   id BIGSERIAL PRIMARY KEY,
   city_id BIGINT NOT NULL REFERENCES cities(id),
   name VARCHAR NOT NULL,
   location GEOGRAPHY(Point, 4326),
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);


CREATE TABLE sub_districts (
   id BIGSERIAL PRIMARY KEY,
   district_id BIGINT NOT NULL REFERENCES districts(id),
   name VARCHAR NOT NULL,
   postal_codes VARCHAR NOT NULL,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);


create table user_addresses (
   id bigserial primary key,
   user_id bigint not null references users(id),
   name varchar NOT NULL,
   phone_number varchar NOT NULL,
   address varchar not null,
   province_id BIGINT NOT NULL REFERENCES provinces(id),
   province varchar NOT NULL,
   city_id BIGINT NOT NULL REFERENCES cities(id),
   city varchar NOT NULL,
   district_id BIGINT NOT NULL REFERENCES districts(id),
   district varchar NOT NULL,
   sub_district_id BIGINT NOT NULL REFERENCES sub_districts(id),
   sub_district varchar NOT NULL,
   postal_code varchar NOT NULL,
   location GEOGRAPHY(Point, 4326) NOT NULL,
   is_active bool not null default false,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);

create table pharmacies (
   id bigserial primary key,
   logo varchar NOT NULL,
   partner_id bigint not null references partners(id),
   name varchar not null,
   address varchar not null,
   province_id BIGINT NOT NULL REFERENCES provinces(id),
   province varchar NOT NULL,
   city_id BIGINT NOT NULL REFERENCES cities(id),
   city varchar NOT NULL,
   district_id BIGINT NOT NULL REFERENCES districts(id),
   district varchar NOT NULL,
   sub_district_id BIGINT NOT NULL REFERENCES sub_districts(id),
   sub_district varchar NOT NULL,
   postal_code bigint NOT NULL,
   location GEOGRAPHY(Point, 4326) NOT NULL,
   is_active bool NOT NULL DEFAULT false,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);

CREATE TABLE logistic_partners (
   id bigserial primary key,
   name VARCHAR NOT NULL,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);


CREATE TABLE pharmacies_logistic_partners(
   id bigserial primary key,
   pharmacy_id bigint NOT NULL REFERENCES pharmacies(id),
   logistic_partner_id bigint NOT NULL REFERENCES logistic_partners(id),
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);


create table pharmacist_details (
   id bigserial primary key,
   pharmacist_id bigint not null references users(id),
   pharmacy_id bigint references pharmacies(id),
   sipa_number varchar not null,
   phone_number varchar not null,
   year_of_experience int not null,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);


create table orders (
   id bigserial primary key,
   user_id bigint not null references users(id),
   total_price decimal(14,2) not null,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);


create table order_details (
   id bigserial primary key,
   order_id bigint not null references orders(id),
   pharmacy_id bigint not null references pharmacies(id),
   logistic_price decimal(14,2) not null,
   status varchar not null,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);


create table products (
	id bigserial primary key,
	product_classification_id bigint not null references product_classifications(id),
	product_form_id bigint null references product_forms(id),
	name varchar not null,
	generic_name varchar not null,
	manufacture varchar not null,
	description varchar not null,
	image varchar[] not null,
	unit_in_pack int null,
	weight decimal(14,2) not null,
	height decimal(14,2) not null,
	length decimal(14,2) not null,
	width decimal(14,2) not null,
	is_active bool not null,
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp,
	deleted_at timestamp null
);


create table product_multi_categories (
	id bigserial primary key,
	product_id bigint not null references products(id),
	product_category_id bigint not null references product_categories(id),
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp,
	deleted_at timestamp null 
);


create table pharmacy_products (
   id bigserial primary key,
   pharmacy_id bigint not null references pharmacies(id),
   product_id bigint not null references products(id),
   stock int not null,
   price decimal(14,2) not null,
   is_active boolean not null,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);


create table order_product_details (
   id bigserial primary key,
   order_detail_id bigint not null references order_details(id),
   pharmacy_product_id bigint not null references pharmacy_products(id),
   quantity int not null,
   price decimal(14,2) not null,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);


create table reset_password_tokens (
   id bigserial primary key,
   user_id bigint not null references users(id),
   token varchar not null,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);


create table carts (
	id bigserial primary key, 
	user_id bigint not null references users(id),
	pharmacy_product_id bigint not null references pharmacy_products(id),
	quantity int not null,
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp,
	deleted_at timestamp null
);


create table verify_email_tokens (
   id bigserial primary key,
   user_id bigint not null references users(id),
   token varchar not null,
   created_at timestamp not null default current_timestamp,
   updated_at timestamp not null default current_timestamp,
   deleted_at timestamp null
);

insert into users (name, email, password, profile_photo, role, is_verified)
values
   ('admin', 'admin@gmail.com', '$2y$10$piFXfOkSHLCrEwHjsAdwEuNCFIoVInLz/XNgXZr7TqYHJm5EPtfbO', 'https://static.thenounproject.com/png/363639-200.png', 'admin', true),
   ('imin', 'imin@gmail.com', '$2y$10$3LLgkx40s1fvO2IA5TeCkuwOOe5wChWwJ.ged4wxSuR5T.PgRCarm', 'https://static.thenounproject.com/png/363639-200.png', 'user', true),
   ('febi', 'febiijuga@gmail.com', '$2y$10$3W77S/wsIvghBVMGKdnt8etsSxEtCkBUA83dTRlLd5Z4jQjxqi9.G', 'https://static.thenounproject.com/png/363639-200.png', 'user', true),
   ('kwang', 'kwang@gmail.com', '$2y$10$2KuXgVQW5gTpb8GIq4.OieZza0xaUBEspHUFtGYnRoiDEqfcQrbuy', 'https://static.thenounproject.com/png/363639-200.png', 'pharmacist', true),
   ('akhdan', 'akhdan@gmail.com', '$2y$10$jnXO5zevtmMeHGAay2XoD.ReJQpuuzefggTyXwTKxVa0SbhZknVXm', 'https://static.thenounproject.com/png/363639-200.png', 'pharmacist', true),
   ('naufal', 'naufal@gmail.com', '$2y$10$Go0c6rpAhx.mICzFjfxWs.JAEftpmoO6ChzCp1rxwFkYBq6qqjjqK', 'https://static.thenounproject.com/png/363639-200.png', 'user', true);


insert into partners (name, year_founded, active_days, start_hour, end_hour, is_active)
values
  ('Sehat Group A', '1998', 'monday,tuesday,wednesday,thursday,friday', '08:00', '23:00', True),
  ('Sehat Group B', '2000', 'monday,tuesday,wednesday,thursday,friday', '09:00', '23:00', True),
  ('Sehat Group C', '2002', 'monday,tuesday,wednesday,thursday,friday', '10:00', '23:00', True),
  ('MedCorp Alpha', '1998', 'monday,tuesday,wednesday,thursday,friday', '07:00', '23:00', True),
  ('MedCorp Beta', '2000', 'monday,tuesday,wednesday,thursday,friday', '10:00', '23:00', True),
  ('Fit Group', '2002', 'monday,tuesday,wednesday,thursday,friday', '11:00', '23:00', True),
  ('NovoCorp Alpha', '2002', 'monday,tuesday,wednesday,thursday,friday', '06:00', '23:00', True),
  ('NovoCorp Beta', '2002', 'monday,tuesday,wednesday,thursday,friday', '09:00', '23:00', True),
  ('Delta Health', '2002', 'monday,tuesday,wednesday,thursday,friday', '08:30', '23:00', True),
  ('MedCare Group', '2002', 'monday,tuesday,wednesday,thursday,friday', '07:30', '23:00', True),
  ('HealthPlus', '2003', 'monday,tuesday,wednesday,thursday,friday', '09:00', '23:00', True),
  ('CareLine', '2005', 'monday,tuesday,wednesday,thursday,friday', '09:00', '23:00', True),
  ('VitalCare', '2000', 'monday,tuesday,wednesday,thursday,friday', '08:30', '23:00', True),
  ('Wellness Group', '2001', 'monday,tuesday,wednesday,thursday,friday', '10:30', '23:00', True),
  ('Happy Health', '2004', 'monday,tuesday,wednesday,thursday,friday', '11:00', '23:00', True);


insert into product_categories (name) values
   ('Pain Relief'),
   ('Antibiotics'),
   ('Vitamins & Supplements'),
   ('Cold & Flu'),
   ('Digestive Health'),
   ('Personal Care'),
   ('Skin Care'),
   ('Hair Care');


insert into product_classifications (name) values
   ('Over the Counter Drugs'),
   ('Prescription Drugs'),
   ('Limited Over the Counter'),
   ('Non Drug');


COPY product_forms(name)
FROM '/data/products/product_forms.csv' CSV HEADER;

COPY products(product_form_id, product_classification_id, name, generic_name, manufacture, description, image, unit_in_pack, weight, height, length, width, is_active)
FROM '/data/products/products.csv' CSV HEADER;


-- insert into products (product_classification_id, product_form_id, name, generic_name, manufacture, description, image, unit_in_pack, weight, height, length, width, is_active) values
--	(1, 1, 'Paracetamol', 'Acetaminophen', 'PharmaCorp', 'Pain relief medication for headaches and fever', array['paracetamol.jpg'], 10, 50.0, 5.0, 10.0, 5.0, true),
--	(2, 2, 'Amoxicillin', 'Amoxicillin', 'MediPharm', 'Antibiotic used for bacterial infections', array['amoxicillin.jpg'], 30, 100.0, 4.0, 10.0, 5.0, true),
--	(3, 3, 'Vitamin C Tablets', 'Ascorbic Acid', 'NutraCorp', 'Vitamin C supplement for immune support', array['vitamin_c.jpg'], 60, 30.0, 3.0, 8.0, 3.0, true),
--	(4, 4, 'Cough Syrup', 'Dextromethorphan', 'WellnessPharma', 'For cough and sore throat relief', array['cough_syrup.jpg'], 200, 150.0, 5.0, 12.0, 6.0, true),
--	(1, 5, 'Antacid Tablets', 'Calcium Carbonate', 'HealthCo', 'For heartburn and indigestion relief', array['antacid_tablets.jpg'], 20, 200.0, 4.0, 9.0, 4.0, false),
--	(2, 6, 'Aloe Vera Gel', 'Aloe Vera', 'SkinCare Inc.', 'Soothing gel for skin irritation and burns', array['aloe_vera.jpg'], 100, 75.0, 7.0, 14.0, 5.0, true),
--	(3, 2, 'Hair Growth Serum', 'Minoxidil', 'BeautyPharm', 'Serum for stimulating hair growth', array['hair_growth.jpg'], 50, 80.0, 6.0, 10.0, 4.0, true);


insert into product_multi_categories (product_id, product_category_id) values
	(1, 1),
	(1, 2),
	(2, 3),
	(3, 3),
	(4, 4),
	(5, 5),
	(6, 6),
	(7, 2);

COPY provinces(id, name)
FROM '/data/provinces.csv' CSV HEADER;

COPY cities(id, province_id, name, location)
FROM '/data/cities.csv' CSV HEADER;

COPY districts(id, city_id, name, location)
FROM '/data/districts.csv' CSV HEADER;

COPY sub_districts(id, district_id, name, postal_codes)
FROM '/data/sub_districts.csv' CSV HEADER;

COPY pharmacies(name, sub_district, district, city, province, location, province_id, city_id, district_id, sub_district_id, postal_code, address, partner_id, logo, is_active)
FROM '/data/pharmacies/pharmacies_jabodetabek.csv' CSV HEADER;

--insert into pharmacy_products (pharmacy_id, product_id, stock, price, is_active) values
--   (1, 1, 100, 5.99, true),
--   (1, 4, 20, 2.00, true),
--   (1, 2, 50, 12.99, true),
--   (2, 3, 200, 8.49, true),
--   (2, 4, 75, 5.99, true),
--   (3, 5, 150, 10.99, true),
--   (3, 6, 30, 7.49, true),
--   (4, 7, 80, 15.00, false),
--   (4, 3, 150, 1.50, true);


COPY pharmacy_products (pharmacy_id, product_id, stock, price, is_active)
FROM '/data/pharmacy_products/pharmacy_products.csv' CSV HEADER;

insert into carts (user_id, pharmacy_product_id, quantity) values
	(2, 1, 1),
	(2, 4, 5),
	(3, 4, 10),
	(3, 2, 5);


INSERT INTO pharmacist_details (pharmacist_id,pharmacy_id,sipa_number,phone_number,year_of_experience) VALUES
(4,5,'111/2122/223','08999123456',5);


insert into orders (user_id, total_price) values
   (2, 75.50),
   (3, 120.00),
   (2, 80.00),
   (2, 110.50),
   (3, 100.00),
   (3, 200.00),
   (3, 200.00);


insert into order_details (order_id, pharmacy_id, logistic_price, status) values
   (1, 1, 10.00, 'Pending'),
   (2, 2, 7.50, 'Shipped'),
   (3, 3, 5.00, 'Delivered'),
   (4, 4, 9.00, 'Cancelled'),
   (5,5,5000,'Pending'),
   (6,5,2000,'Shipped'),
   (7,5,5000,'Processing');


insert into order_product_details (order_detail_id, pharmacy_product_id, quantity, price) values
   (1, 1, 5, 5.99),
   (1, 2, 3, 12.99),
   (2, 3, 2, 8.49),
   (2, 4, 1, 5.99),
   (5, 5, 1, 10.99),
   (6, 6, 2, 7.49),
   (7, 7, 3, 15.00);

INSERT INTO logistics (name, price)
VALUES
('same day',1000),
('instant',2500);

INSERT INTO public.user_addresses (user_id, "name", phone_number, address, province_id, province, city_id, city, district_id, district, sub_district_id, sub_district, postal_code, "location", is_active, created_at, updated_at, deleted_at)
VALUES
(2, 'Jonathan', '089505123456', 'Jalan Pos', 31, 'DKI JAKARTA', 3174, 'JAKARTA BARAT', 3174010, 'KEMBANGAN', 3174010002, 'SRENGSENG', '11630', 'SRID=4326;POINT (106.74130492346 -6.191140471555)'::public.geography, true, '2025-01-27 18:26:57.025', '2025-01-27 18:27:02.874', NULL);
