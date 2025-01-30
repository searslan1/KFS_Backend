-- users tablosu
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    user_type VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- auth_users tablosu
CREATE TABLE auth_users (
    user_id INTEGER PRIMARY KEY REFERENCES users(user_id),
    email_verified BOOLEAN DEFAULT FALSE,
    phone_verified BOOLEAN DEFAULT FALSE,
    password_reset_token VARCHAR(255),
    password_reset_expires TIMESTAMP,
    failed_login_attempts INTEGER DEFAULT 0,
    account_locked_until TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- usertype_profile tablosu
CREATE TABLE usertype_profile (
    user_id INTEGER PRIMARY KEY REFERENCES users(user_id),
    name VARCHAR(255) NOT NULL,
    id VARCHAR(50) UNIQUE NOT NULL,
    registration_number VARCHAR(50) UNIQUE,
    tax_office VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20)
);

-- company_representatives tablosu
CREATE TABLE company_representatives (
    rep_id SERIAL PRIMARY KEY,
    company_id INTEGER REFERENCES usertype_profile(user_id),
    user_id INTEGER REFERENCES users(user_id),
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- user_profiles tablosu
CREATE TABLE user_profiles (
    profile_id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users(user_id),
    profile_picture_url TEXT,
    bio TEXT,
    website_url TEXT,
    social_media JSONB,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- addresses tablosu
CREATE TABLE addresses (
    address_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    country VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    district VARCHAR(100),
    street VARCHAR(255),
    zip_code VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- bank_accounts tablosu
CREATE TABLE bank_accounts (
    account_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    bank_name VARCHAR(255) NOT NULL,
    account_number VARCHAR(50) UNIQUE NOT NULL,
    iban VARCHAR(50) UNIQUE NOT NULL,
    account_type VARCHAR(20) NOT NULL
);

-- media_files tablosu
CREATE TABLE media_files (
    media_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER REFERENCES users(user_id),
    file_url TEXT NOT NULL,
    file_type VARCHAR(20) NOT NULL,
    file_size INTEGER NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- user_sessions tablosu
CREATE TABLE user_sessions (
    session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER REFERENCES users(user_id),
    ip_address INET NOT NULL,
    user_agent TEXT,
    device_info TEXT,
    login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    logout_time TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    last_activity TIMESTAMP,
    CONSTRAINT idx_user_sessions_user_id UNIQUE (user_id),
    CONSTRAINT idx_user_sessions_login_time UNIQUE (login_time)
);

-- user_security tablosu
CREATE TABLE user_security (
    user_id INTEGER PRIMARY KEY REFERENCES users(user_id),
    has_mfa_enabled BOOLEAN DEFAULT FALSE,
    mfa_secret VARCHAR(255),
    recovery_codes JSONB
);

-- user_activity_logs tablosu
CREATE TABLE user_activity_logs (
    log_id BIGSERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    action VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address INET,
    user_agent TEXT,
    archived BOOLEAN DEFAULT FALSE
);

-- security_alerts tablosu
CREATE TABLE security_alerts (
    alert_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    alert_type VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- user_notifications tablosu
CREATE TABLE user_notifications (
    notification_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- investor_profiles tablosu
CREATE TABLE investor_profiles (
    investor_id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users(user_id),
    investment_focus TEXT,
    min_investment_amount DECIMAL(15,2),
    max_investment_amount DECIMAL(15,2),
    past_investments TEXT,
    risk_profile VARCHAR(50),
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- entrepreneur_profiles tablosu
CREATE TABLE entrepreneur_profiles (
    entrepreneur_id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users(user_id),
    startup_name VARCHAR(255) NOT NULL,
    industry VARCHAR(100) NOT NULL,
    funding_needed DECIMAL(15,2) NOT NULL,
    business_model TEXT,
    pitch_deck_url TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- campaign_profile tablosu
CREATE TABLE campaign_profile (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    campaign_logo TEXT,
    entrepreneur_name VARCHAR(255),
    campaign_name VARCHAR(255) NOT NULL,
    campaign_description TEXT,
    about_project TEXT,
    campaign_summary TEXT,
    goal_coverage_subject TEXT,
    entrepreneur_stage_id INTEGER,
    location VARCHAR(255),
    category VARCHAR(100),
    business_models_id INTEGER,
    sector VARCHAR(100),
    entrepreneurs_mails TEXT,
    is_past_campaign BOOLEAN DEFAULT FALSE,
    campaign_status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- prize tablosu
CREATE TABLE prize (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaign_profile(id),
    prize_date DATE,
    prize_description TEXT,
    prize_path TEXT,
    awarding_organization VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- patent tablosu
CREATE TABLE patent (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaign_profile(id),
    document_no VARCHAR(50),
    description TEXT,
    document_path TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- law tablosu
CREATE TABLE law (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaign_profile(id),
    permission_subject VARCHAR(255),
    permission_path TEXT,
    permission_description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- team_members tablosu
CREATE TABLE team_members (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaign_profile(id),
    photograph TEXT,
    members_name VARCHAR(100),
    members_surname VARCHAR(100),
    members_title VARCHAR(100),
    resume TEXT,
    biography TEXT,
    members_task TEXT,
    members_responsibility TEXT,
    entrepreneur_link_member TEXT,
    members_mail VARCHAR(255),
    social_links JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- business_models tablosu
CREATE TABLE business_models (
    id SERIAL PRIMARY KEY,
    model TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- entrepreneur_stage tablosu
CREATE TABLE entrepreneur_stage (
    id SERIAL PRIMARY KEY,
    stage TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- product_production_models tablosu
CREATE TABLE product_production_models (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaign_profile(id),
    product_summary TEXT,
    product_about TEXT,
    product_path TEXT,
    product_problem TEXT,
    product_solution TEXT,
    product_evaluation TEXT,
    products_development_stage_summary TEXT,
    products_development_stage_summary_path TEXT,
    products_production_stage_summary TEXT,
    products_production_stage_summary_path TEXT,
    side_product_summary TEXT,
    side_product_summary_path TEXT,
    analysis_summary TEXT,
    analysis_summary_path TEXT,
    AR_GE_summary TEXT,
    AR_GE_summary_path TEXT,
    preview_sales_summary TEXT,
    preview_sales_summary_path TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- other_subject tablosu
CREATE TABLE other_subject (
    id SERIAL PRIMARY KEY,
    product_production_models_id INTEGER REFERENCES product_production_models(id),
    subject TEXT,
    path TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- financial_table tablosu
CREATE TABLE financial_table (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaign_profile(id),
    investment_year_1 DECIMAL(15,2),
    investment_year_2 DECIMAL(15,2),
    investment_year_3 DECIMAL(15,2),
    investment_year_4 DECIMAL(15,2),
    investment_year_5 DECIMAL(15,2),
    total DECIMAL(15,2),
    profit_estimates_text TEXT,
    explanation_text TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- product_revenue tablosu
CREATE TABLE product_revenue (
    id SERIAL PRIMARY KEY,
    financial_table_id INTEGER REFERENCES financial_table(id),
    product_name VARCHAR(255),
    sales_price DECIMAL(15,2),
    direct_cost DECIMAL(15,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- sales_targets tablosu
CREATE TABLE sales_targets (
    id SERIAL PRIMARY KEY,
    product_revenue_id INTEGER REFERENCES product_revenue(id),
    target_type VARCHAR(100),
    year_1 DECIMAL(15,2),
    year_2 DECIMAL(15,2),
    year_3 DECIMAL(15,2),
    year_4 DECIMAL(15,2),
    year_5 DECIMAL(15,2),
    total DECIMAL(15,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- budget tablosu
CREATE TABLE budget (
    id SERIAL PRIMARY KEY,
    financial_table_id INTEGER REFERENCES financial_table(id),
    category VARCHAR(100),
    sub_category VARCHAR(100),
    year_1 DECIMAL(15,2),
    year_2 DECIMAL(15,2),
    year_3 DECIMAL(15,2),
    year_4 DECIMAL(15,2),
    year_5 DECIMAL(15,2),
    total DECIMAL(15,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- document tablosu
CREATE TABLE document (
    id SERIAL PRIMARY KEY,
    financial_table_id INTEGER REFERENCES financial_table(id),
    document_name VARCHAR(255),
    document_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- funding tablosu
CREATE TABLE funding (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaign_profile(id),
    interference_value DECIMAL(15,2),
    use_time TEXT,
    evaluation_report TEXT,
    need_amount_fund DECIMAL(15,2),
    amount_given_share DECIMAL(15,2),
    number_sale_share INTEGER,
    post_funding_capital DECIMAL(15,2),
    sales_share_price DECIMAL(15,2),
    sales_nominal_share_price DECIMAL(15,2),
    place_funds_collected_id INTEGER REFERENCES place_funds_collected(id),
    additional_sources_financing_id INTEGER REFERENCES additional_sources_financing(id),
    comparison_current_postfunding TEXT,
    basic_information TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- place_funds_collected tablosu
CREATE TABLE place_funds_collected (
    id SERIAL PRIMARY KEY,
    usage_start_time TIMESTAMP,
    usage_finish_time TIMESTAMP,
    description TEXT,
    amount DECIMAL(15,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- additional_sources_financing tablosu
CREATE TABLE additional_sources_financing (
    id SERIAL PRIMARY KEY,
    provision_time TIMESTAMP,
    description TEXT,
    amount DECIMAL(15,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- establishment tablosu
CREATE TABLE establishment (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaign_profile(id),
    title VARCHAR(255),
    capitalization DECIMAL(15,2),
    city VARCHAR(100),
    district VARCHAR(100),
    address TEXT,
    postfunding_partners_id INTEGER REFERENCES postfunding_partners(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- postfunding_partners tablosu
CREATE TABLE postfunding_partners (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    task TEXT,
    university VARCHAR(255),
    average_not DECIMAL(3,2),
    resume TEXT,
    citizenship VARCHAR(50),
    share_in_capital_amount DECIMAL(15,2),
    share_in_capital_rate DECIMAL(5,2),
    vote BOOLEAN,
    concession TEXT,
    campaign_relation TEXT,
    work_experience TEXT,
    areas_of_specialization TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- visual_video tablosu
CREATE TABLE visual_video (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaign_profile(id),
    display_photograph TEXT,
    other_photograph TEXT,
    video_link TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- other_document tablosu
CREATE TABLE other_document (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaign_profile(id),
    file_path TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- risk tablosu
CREATE TABLE risk (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaign_profile(id),
    risk_type VARCHAR(50),
    risk_description TEXT,
    risk_mitigation TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- campaign_status_logs tablosu
CREATE TABLE campaign_status_logs (
    log_id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaign_profile(id),
    previous_status VARCHAR(50),
    new_status VARCHAR(50),
    changed_by INTEGER REFERENCES users(user_id),
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- campaign_investments tablosu
CREATE TABLE campaign_investments (
    investment_id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaign_profile(id),
    investor_id INTEGER REFERENCES users(user_id),
    amount DECIMAL(15,2) NOT NULL,
    investment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending'
);

-- transactions tablosu
CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    campaign_id INTEGER REFERENCES campaign_profile(id) DEFAULT NULL,
    amount DECIMAL(15,2) NOT NULL,
    transaction_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    reference_code VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- rate_limits tablosu
CREATE TABLE rate_limits (
    limit_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    endpoint VARCHAR(255),
    request_count INTEGER DEFAULT 0,
    last_request TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- kyc_verifications tablosu
CREATE TABLE kyc_verifications (
    kyc_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    status VARCHAR(20) DEFAULT 'pending',
    verification_details JSONB,
    verified_at TIMESTAMP
);

-- aml_logs tablosu
CREATE TABLE aml_logs (
    aml_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    transaction_id INTEGER REFERENCES transactions(transaction_id),
    aml_status VARCHAR(20) DEFAULT 'pending',
    flagged_reason TEXT,
    checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- user_permissions_log tablosu
CREATE TABLE user_permissions_log (
    log_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    changed_by INTEGER REFERENCES users(user_id),
    previous_permissions TEXT,
    new_permissions TEXT,
    change_reason TEXT,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- investor_portfolio tablosu
CREATE TABLE investor_portfolio (
    portfolio_id SERIAL PRIMARY KEY,
    investor_id INTEGER REFERENCES users(user_id),
    total_investment DECIMAL(15,2) DEFAULT 0,
    active_investments JSONB,
    total_shares INTEGER DEFAULT 0,
    last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- user_agreements tablosu
CREATE TABLE user_agreements (
    agreement_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    agreement_type VARCHAR(100),
    agreement_version VARCHAR(10),
    agreed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- user_interactions tablosu
CREATE TABLE user_interactions (
    interaction_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    campaign_id INTEGER REFERENCES campaign_profile(id),
    interaction_type VARCHAR(50),
    interaction_data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- payment_details tablosu
CREATE TABLE payment_details (
    payment_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    transaction_id INTEGER REFERENCES transactions(transaction_id),
    payment_provider VARCHAR(50),
    payment_status VARCHAR(20) DEFAULT 'pending',
    refund_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- campaign_analytics tablosu
CREATE TABLE campaign_analytics (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaign_profile(id),
    total_visits INTEGER DEFAULT 0,
    total_watch_time DECIMAL(10,2),
    unique_visitors INTEGER DEFAULT 0,
    total_investments INTEGER DEFAULT 0,
    total_investment_amount DECIMAL(15,2) DEFAULT 0,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- user_feedback tablosu
CREATE TABLE user_feedback (
    feedback_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    campaign_id INTEGER REFERENCES campaign_profile(id) DEFAULT NULL,
    feedback_type VARCHAR(50),
    feedback_text TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- user_activity_tracking tablosu
CREATE TABLE user_activity_tracking (
    tracking_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    page_url TEXT NOT NULL,
    session_id UUID REFERENCES user_sessions(session_id),
    time_spent DECIMAL(10,2),
    event_type VARCHAR(50),
    event_data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- system_notifications tablosu
CREATE TABLE system_notifications (
    notification_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    notification_type VARCHAR(50),
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- investment_recommendations tablosu
CREATE TABLE investment_recommendations (
    id SERIAL PRIMARY KEY,
    investor_id INTEGER REFERENCES users(user_id),
    recommended_campaigns JSONB,
    recommendation_reason TEXT,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- investment_performance tablosu
CREATE TABLE investment_performance (
    id SERIAL PRIMARY KEY,
    investor_id INTEGER REFERENCES users(user_id),
    total_invested DECIMAL(15,2) DEFAULT 0,
    total_return DECIMAL(15,2) DEFAULT 0,
    current_value DECIMAL(15,2),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
