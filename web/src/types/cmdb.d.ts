export interface Provider {
  id: number;
  name: string;
  type: 'aliyun' | 'tencent' | 'aws';
  access_key: string;
  secret_key: string;
  region: string;
  description: string;
  status: 'enabled' | 'disabled';
  created_at: string;
  updated_at: string;
}