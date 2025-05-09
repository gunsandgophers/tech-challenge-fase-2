CREATE TABLE orders
(
    id uuid NOT NULL,
    customer_id uuid,
    items JSON[] DEFAULT ARRAY[]::JSON[],
    payment_status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    preparation_status VARCHAR(20) NOT NULL DEFAULT 'AWAITING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_order_id PRIMARY KEY (id),
    CONSTRAINT payment_status CHECK (payment_status IN ('PENDING', 'PAID', 'REJECTED', 'AWAITING_PAYMENT')),
    CONSTRAINT preparation_status CHECK (preparation_status IN ('AWAITING', 'RECEIVED', 'IN_PREPARATION', 'READY', 'FINISHED'))
);
