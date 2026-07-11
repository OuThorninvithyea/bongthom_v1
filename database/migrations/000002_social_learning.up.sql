BEGIN;

-- This migration depends on the existing public.tbl_users table.
-- MVP assumptions:
--   * Following is one-way.
--   * A course has one instructor.
--   * Lessons support text, video, and file content.
--   * Course pricing is metadata only; payments are a separate module.
--   * Messaging, notifications, and quizzes are outside this migration.

CREATE TABLE tbl_media_assets (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL REFERENCES tbl_users(id) ON DELETE RESTRICT,
    media_type VARCHAR(16) NOT NULL
        CHECK (media_type IN ('image', 'video', 'file')),
    storage_provider VARCHAR(32) NOT NULL DEFAULT 'local',
    bucket_name TEXT NOT NULL CHECK (LENGTH(TRIM(bucket_name)) > 0),
    object_key TEXT NOT NULL CHECK (LENGTH(TRIM(object_key)) > 0),
    mime_type VARCHAR(128) NOT NULL CHECK (LENGTH(TRIM(mime_type)) > 0),
    size_bytes BIGINT NOT NULL CHECK (size_bytes >= 0),
    width INTEGER CHECK (width IS NULL OR width > 0),
    height INTEGER CHECK (height IS NULL OR height > 0),
    duration_ms BIGINT CHECK (duration_ms IS NULL OR duration_ms >= 0),
    status VARCHAR(16) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'ready', 'failed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (storage_provider, bucket_name, object_key)
);

CREATE TABLE tbl_user_profiles (
    user_id BIGINT PRIMARY KEY REFERENCES tbl_users(id) ON DELETE CASCADE,
    display_name VARCHAR(120) NOT NULL
        CHECK (LENGTH(TRIM(display_name)) BETWEEN 1 AND 120),
    bio TEXT CHECK (bio IS NULL OR LENGTH(bio) <= 2000),
    headline VARCHAR(180),
    avatar_media_id BIGINT REFERENCES tbl_media_assets(id) ON DELETE SET NULL,
    cover_media_id BIGINT REFERENCES tbl_media_assets(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tbl_user_follows (
    follower_id BIGINT NOT NULL REFERENCES tbl_users(id) ON DELETE CASCADE,
    following_id BIGINT NOT NULL REFERENCES tbl_users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id),
    CHECK (follower_id <> following_id)
);

CREATE TABLE tbl_courses (
    id BIGSERIAL PRIMARY KEY,
    instructor_id BIGINT NOT NULL REFERENCES tbl_users(id) ON DELETE RESTRICT,
    title VARCHAR(200) NOT NULL
        CHECK (LENGTH(TRIM(title)) BETWEEN 1 AND 200),
    slug VARCHAR(220) NOT NULL
        CHECK (slug ~ '^[a-z0-9]+(?:-[a-z0-9]+)*$'),
    description TEXT,
    thumbnail_media_id BIGINT REFERENCES tbl_media_assets(id) ON DELETE SET NULL,
    level VARCHAR(16) NOT NULL DEFAULT 'beginner'
        CHECK (level IN ('beginner', 'intermediate', 'advanced')),
    visibility VARCHAR(16) NOT NULL DEFAULT 'public'
        CHECK (visibility IN ('public', 'private', 'unlisted')),
    status VARCHAR(16) NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'published', 'archived')),
    price_minor BIGINT NOT NULL DEFAULT 0 CHECK (price_minor >= 0),
    currency_code CHAR(3)
        CHECK (currency_code IS NULL OR currency_code ~ '^[A-Z]{3}$'),
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CHECK (price_minor = 0 OR currency_code IS NOT NULL),
    CHECK (status <> 'published' OR published_at IS NOT NULL)
);

CREATE TABLE tbl_course_sections (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL REFERENCES tbl_courses(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL
        CHECK (LENGTH(TRIM(title)) BETWEEN 1 AND 200),
    position INTEGER NOT NULL CHECK (position >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_id, position),
    UNIQUE (id, course_id)
);

CREATE TABLE tbl_lessons (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL,
    section_id BIGINT NOT NULL,
    title VARCHAR(200) NOT NULL
        CHECK (LENGTH(TRIM(title)) BETWEEN 1 AND 200),
    description TEXT,
    lesson_type VARCHAR(16) NOT NULL
        CHECK (lesson_type IN ('video', 'text', 'file')),
    media_id BIGINT REFERENCES tbl_media_assets(id) ON DELETE SET NULL,
    text_content TEXT,
    duration_seconds INTEGER
        CHECK (duration_seconds IS NULL OR duration_seconds >= 0),
    position INTEGER NOT NULL CHECK (position >= 0),
    is_preview BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(16) NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'published', 'archived')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (section_id, course_id)
        REFERENCES tbl_course_sections(id, course_id) ON DELETE CASCADE,
    UNIQUE (section_id, position),
    UNIQUE (id, course_id),
    CHECK (
        (lesson_type = 'text' AND text_content IS NOT NULL
            AND LENGTH(TRIM(text_content)) > 0)
        OR
        (lesson_type IN ('video', 'file') AND media_id IS NOT NULL)
    )
);

CREATE TABLE tbl_enrollments (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL REFERENCES tbl_courses(id) ON DELETE RESTRICT,
    user_id BIGINT NOT NULL REFERENCES tbl_users(id) ON DELETE RESTRICT,
    status VARCHAR(16) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'completed', 'cancelled')),
    enrolled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_id, user_id),
    UNIQUE (id, course_id),
    CHECK (
        (status = 'completed' AND completed_at IS NOT NULL)
        OR
        (status <> 'completed' AND completed_at IS NULL)
    )
);

CREATE TABLE tbl_lesson_progress (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL,
    enrollment_id BIGINT NOT NULL,
    lesson_id BIGINT NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'not_started'
        CHECK (status IN ('not_started', 'in_progress', 'completed')),
    progress_percent NUMERIC(5,2) NOT NULL DEFAULT 0
        CHECK (progress_percent BETWEEN 0 AND 100),
    last_position_seconds INTEGER NOT NULL DEFAULT 0
        CHECK (last_position_seconds >= 0),
    completed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (enrollment_id, course_id)
        REFERENCES tbl_enrollments(id, course_id) ON DELETE CASCADE,
    FOREIGN KEY (lesson_id, course_id)
        REFERENCES tbl_lessons(id, course_id) ON DELETE CASCADE,
    UNIQUE (enrollment_id, lesson_id),
    CHECK (
        (status = 'completed' AND progress_percent = 100
            AND completed_at IS NOT NULL)
        OR
        (status <> 'completed' AND progress_percent < 100
            AND completed_at IS NULL)
    ),
    CHECK (status <> 'not_started' OR progress_percent = 0)
);

CREATE TABLE tbl_posts (
    id BIGSERIAL PRIMARY KEY,
    author_id BIGINT NOT NULL REFERENCES tbl_users(id) ON DELETE RESTRICT,
    course_id BIGINT REFERENCES tbl_courses(id) ON DELETE SET NULL,
    shared_post_id BIGINT REFERENCES tbl_posts(id) ON DELETE SET NULL,
    body TEXT
        CHECK (body IS NULL OR LENGTH(TRIM(body)) BETWEEN 1 AND 10000),
    visibility VARCHAR(16) NOT NULL DEFAULT 'public'
        CHECK (visibility IN ('public', 'followers', 'course', 'private')),
    status VARCHAR(16) NOT NULL DEFAULT 'published'
        CHECK (status IN ('draft', 'published', 'hidden')),
    comments_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CHECK (visibility <> 'course' OR course_id IS NOT NULL),
    CHECK (shared_post_id IS NULL OR shared_post_id <> id)
);

CREATE TABLE tbl_post_media (
    post_id BIGINT NOT NULL REFERENCES tbl_posts(id) ON DELETE CASCADE,
    media_id BIGINT NOT NULL REFERENCES tbl_media_assets(id) ON DELETE RESTRICT,
    position INTEGER NOT NULL DEFAULT 0 CHECK (position >= 0),
    PRIMARY KEY (post_id, media_id),
    UNIQUE (post_id, position)
);

CREATE TABLE tbl_comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES tbl_posts(id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL REFERENCES tbl_users(id) ON DELETE RESTRICT,
    parent_comment_id BIGINT,
    body TEXT NOT NULL
        CHECK (LENGTH(TRIM(body)) BETWEEN 1 AND 5000),
    status VARCHAR(16) NOT NULL DEFAULT 'published'
        CHECK (status IN ('published', 'hidden')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (id, post_id),
    FOREIGN KEY (parent_comment_id, post_id)
        REFERENCES tbl_comments(id, post_id) ON DELETE CASCADE
);

CREATE TABLE ref_reaction_types (
    id SMALLSERIAL PRIMARY KEY,
    code VARCHAR(24) NOT NULL UNIQUE,
    label VARCHAR(50) NOT NULL,
    position SMALLINT NOT NULL DEFAULT 0 CHECK (position >= 0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

INSERT INTO ref_reaction_types (code, label, position)
VALUES
    ('like', 'Like', 1),
    ('love', 'Love', 2),
    ('haha', 'Haha', 3),
    ('wow', 'Wow', 4),
    ('sad', 'Sad', 5),
    ('angry', 'Angry', 6);

CREATE TABLE tbl_post_reactions (
    post_id BIGINT NOT NULL REFERENCES tbl_posts(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES tbl_users(id) ON DELETE CASCADE,
    reaction_type_id SMALLINT NOT NULL
        REFERENCES ref_reaction_types(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (post_id, user_id)
);

CREATE TABLE tbl_comment_reactions (
    comment_id BIGINT NOT NULL REFERENCES tbl_comments(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES tbl_users(id) ON DELETE CASCADE,
    reaction_type_id SMALLINT NOT NULL
        REFERENCES ref_reaction_types(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (comment_id, user_id)
);

CREATE TABLE tbl_stories (
    id BIGSERIAL PRIMARY KEY,
    author_id BIGINT NOT NULL REFERENCES tbl_users(id) ON DELETE RESTRICT,
    media_id BIGINT NOT NULL REFERENCES tbl_media_assets(id) ON DELETE RESTRICT,
    caption VARCHAR(500),
    visibility VARCHAR(16) NOT NULL DEFAULT 'followers'
        CHECK (visibility IN ('public', 'followers', 'private')),
    status VARCHAR(16) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'hidden')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '24 hours'),
    deleted_at TIMESTAMPTZ,
    CHECK (caption IS NULL OR LENGTH(TRIM(caption)) BETWEEN 1 AND 500),
    CHECK (expires_at > created_at)
);

CREATE TABLE tbl_story_views (
    story_id BIGINT NOT NULL REFERENCES tbl_stories(id) ON DELETE CASCADE,
    viewer_id BIGINT NOT NULL REFERENCES tbl_users(id) ON DELETE CASCADE,
    viewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (story_id, viewer_id)
);

CREATE TABLE tbl_story_reactions (
    story_id BIGINT NOT NULL REFERENCES tbl_stories(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES tbl_users(id) ON DELETE CASCADE,
    reaction_type_id SMALLINT NOT NULL
        REFERENCES ref_reaction_types(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (story_id, user_id)
);

CREATE UNIQUE INDEX uq_courses_slug_active
    ON tbl_courses (slug)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_media_owner
    ON tbl_media_assets (owner_id, created_at DESC);

CREATE INDEX idx_follows_following
    ON tbl_user_follows (following_id, created_at DESC);

CREATE INDEX idx_courses_instructor
    ON tbl_courses (instructor_id, created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_courses_published
    ON tbl_courses (published_at DESC, id DESC)
    WHERE status = 'published' AND deleted_at IS NULL;

CREATE INDEX idx_sections_course
    ON tbl_course_sections (course_id, position);

CREATE INDEX idx_lessons_section
    ON tbl_lessons (section_id, position);

CREATE INDEX idx_enrollments_user
    ON tbl_enrollments (user_id, status);

CREATE INDEX idx_progress_enrollment
    ON tbl_lesson_progress (enrollment_id, updated_at DESC);

CREATE INDEX idx_posts_feed
    ON tbl_posts (created_at DESC, id DESC)
    WHERE status = 'published' AND deleted_at IS NULL;

CREATE INDEX idx_posts_author
    ON tbl_posts (author_id, created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_posts_course
    ON tbl_posts (course_id, created_at DESC)
    WHERE course_id IS NOT NULL AND deleted_at IS NULL;

CREATE INDEX idx_posts_shared
    ON tbl_posts (shared_post_id)
    WHERE shared_post_id IS NOT NULL;

CREATE INDEX idx_comments_post
    ON tbl_comments (post_id, created_at ASC)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_stories_author
    ON tbl_stories (author_id, created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_stories_active
    ON tbl_stories (expires_at, created_at DESC)
    WHERE status = 'active' AND deleted_at IS NULL;

CREATE OR REPLACE FUNCTION set_kaifin_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_media_updated_at
    BEFORE UPDATE ON tbl_media_assets
    FOR EACH ROW EXECUTE FUNCTION set_kaifin_updated_at();

CREATE TRIGGER trg_profiles_updated_at
    BEFORE UPDATE ON tbl_user_profiles
    FOR EACH ROW EXECUTE FUNCTION set_kaifin_updated_at();

CREATE TRIGGER trg_courses_updated_at
    BEFORE UPDATE ON tbl_courses
    FOR EACH ROW EXECUTE FUNCTION set_kaifin_updated_at();

CREATE TRIGGER trg_sections_updated_at
    BEFORE UPDATE ON tbl_course_sections
    FOR EACH ROW EXECUTE FUNCTION set_kaifin_updated_at();

CREATE TRIGGER trg_lessons_updated_at
    BEFORE UPDATE ON tbl_lessons
    FOR EACH ROW EXECUTE FUNCTION set_kaifin_updated_at();

CREATE TRIGGER trg_enrollments_updated_at
    BEFORE UPDATE ON tbl_enrollments
    FOR EACH ROW EXECUTE FUNCTION set_kaifin_updated_at();

CREATE TRIGGER trg_progress_updated_at
    BEFORE UPDATE ON tbl_lesson_progress
    FOR EACH ROW EXECUTE FUNCTION set_kaifin_updated_at();

CREATE TRIGGER trg_posts_updated_at
    BEFORE UPDATE ON tbl_posts
    FOR EACH ROW EXECUTE FUNCTION set_kaifin_updated_at();

CREATE TRIGGER trg_comments_updated_at
    BEFORE UPDATE ON tbl_comments
    FOR EACH ROW EXECUTE FUNCTION set_kaifin_updated_at();

CREATE TRIGGER trg_post_reactions_updated_at
    BEFORE UPDATE ON tbl_post_reactions
    FOR EACH ROW EXECUTE FUNCTION set_kaifin_updated_at();

CREATE TRIGGER trg_comment_reactions_updated_at
    BEFORE UPDATE ON tbl_comment_reactions
    FOR EACH ROW EXECUTE FUNCTION set_kaifin_updated_at();

CREATE TRIGGER trg_stories_updated_at
    BEFORE UPDATE ON tbl_stories
    FOR EACH ROW EXECUTE FUNCTION set_kaifin_updated_at();

CREATE TRIGGER trg_story_reactions_updated_at
    BEFORE UPDATE ON tbl_story_reactions
    FOR EACH ROW EXECUTE FUNCTION set_kaifin_updated_at();

COMMIT;
