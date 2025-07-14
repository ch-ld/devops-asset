// Types for CMDB Provider and Host

export type ProviderType = 'aliyun' | 'tencent' | 'aws' | 'manual';

export interface Provider {
  id: number;
  name: string;
  type: ProviderType;
  access_key: string;
  secret_key: string;
  region: string;
  description: string;
  status: 'enabled' | 'disabled';
  created_at: string;
  updated_at: string;
}

export interface HostGroup {
  id: number;
  name: string;
  description: string;
  parent_id?: number;
  path: string;
  sort: number;
  created_at: string;
  updated_at: string;
}

export interface Host {
  id: number;
  provider_id?: number;
  provider?: Provider;
  instance_id: string;
  name: string;
  resource_type: string;
  region: string;
  username: string;
  password: string;
  public_ip: string[];
  private_ip: string[];
  configuration: {
    cpu_cores?: number;
    memory_size?: number;
    instance_type?: string;
    zone_id?: string;
    vpc_id?: string;
    [key: string]: any;
  };
  os: string;
  status: string;
  expired_at?: string;
  provider_type: string;
  group_id?: number;
  group?: HostGroup;
  tags?: string[];
  extra_fields?: { [key: string]: any };
  remark: string;
  created_at: string;
  updated_at: string;
}

export interface CustomFieldDefinition {
  id: number;
  name: string;
  label: string;
  type: 'text' | 'select' | 'number' | 'textarea' | 'date' | 'checkbox';
  options?: string;
  required: boolean;
  description: string;
  order: number;
}

export interface BatchOperationResult {
  total: number;
  success: number;
  failed: number;
  failed_msg: string[];
}

export interface SSHResult {
  host_id: number;
  success: boolean;
  output: string;
  error: string;
}

export interface SFTPResult {
  host_id: number;
  success: boolean;
  error: string;
}

export interface HostAlert {
  host: Host;
  alert_type: 'expired' | 'expiring' | 'error' | 'abnormal';
  message: string;
} 