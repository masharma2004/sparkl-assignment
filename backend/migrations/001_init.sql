CREATE TABLE users (
    user_id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    role VARCHAR(10) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE questions (
    question_id BIGSERIAL PRIMARY KEY,
    prompt TEXT NOT NULL,
    options JSONB NOT NULL,
    correct_options JSONB NOT NULL,
    solution TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE quizzes (
    quiz_id BIGSERIAL PRIMARY KEY,
    title VARCHAR(250) NOT NULL,
    question_count INT NOT NULL,
    total_marks INT NOT NULL,
    duration_minutes INT NOT NULL,
    created_by BIGINT NOT NULL REFERENCES users(user_id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE quiz_questions (
    quiz_question_id BIGSERIAL PRIMARY KEY,
    quiz_id BIGINT NOT NULL REFERENCES quizzes(quiz_id) ON DELETE CASCADE,
    question_id BIGINT NOT NULL REFERENCES questions(question_id),
    sequence_number INT NOT NULL,
    marks INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE quiz_attempts (
    quiz_attempt_id BIGSERIAL PRIMARY KEY,
    quiz_id BIGINT NOT NULL REFERENCES quizzes(quiz_id) ON DELETE CASCADE,
    student_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL,
    started_at TIMESTAMP NOT NULL,
    submitted_at TIMESTAMP NULL,
    score INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE attempt_answers (
    attempt_answer_id BIGSERIAL PRIMARY KEY,
    attempt_id BIGINT NOT NULL REFERENCES quiz_attempts(quiz_attempt_id) ON DELETE CASCADE,
    quiz_question_id BIGINT NOT NULL REFERENCES quiz_questions(quiz_question_id) ON DELETE CASCADE,
    chosen_options JSONB NOT NULL,
    awarded_marks INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
