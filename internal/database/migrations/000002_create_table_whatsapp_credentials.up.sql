CREATE TABLE whatsapp_credentials (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  phone_number VARCHAR (255) NOT NULL UNIQUE,
  jid VARCHAR (255) NOT NULL UNIQUE,
  is_connected BOOLEAN DEFAULT FALSE,
  session_initiated_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
)