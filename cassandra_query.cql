CREATE KEYSPACE qollege_benchmarks WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1};

CREATE TABLE academic_structures (
  id text,
  full_name text,
  short_name text,
  letter text,
  parent_id text,
  structure_type_id text,
  level_id text,
  entity_id text,
  operator text,
  status int,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (id,structure_type_id,entity_id)
);

CREATE TABLE courses (
  id text,
  full_name text,
  short_name text,
  category_id text,
  entity_id text,
  credits int,
  marks int,
  operator text,
  status int,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  parent_course text,
  PRIMARY KEY (id,entity_id,category_id)
);

CREATE TABLE classes_records (
  id text,
  classe_id text,
  academic_year_id text,
  student_id text,
  operator text,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (id,academic_year_id,classe_id)
);

CREATE TABLE periods (
  id text,
  name text,
  description text,
  type text,
  start_date text,
  end_date text,
  parent_id text,
  status int,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (id, type,parent_id)
);

CREATE TABLE students (
  id text,
  name text,
  email text,
  password text,
  phone text,
  sex text,
  reg_no text,
  dob text,
  place_of_birth text,
  address text,
  country text,
  civil_status text,
  operator text,
  status int,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  default_language text,
  image text,
  PRIMARY KEY (id,reg_no)
);

CREATE TABLE entities (
  id text,
  full_name text,
  short_name text,
  parent_id text,
  entity_type_id text,
  email text,
  phone text,
  country text,
  operator text,
  status int,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (id,entity_type_id,full_name)
);

CREATE TABLE marks (
  id text,
  assessment_id text,
  academic_year_id text,
  term_id text,
  classe_id text,
  course_id text,
  out_of_marks text,
  operator text,
  status int,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (id,assessment_id,academic_year_id,term_id)
);