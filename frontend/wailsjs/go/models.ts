export namespace app {
	
	export class CreateMcpServerRequest {
	    name: string;
	    type: string;
	    command: string;
	    args: string;
	    env: string;
	    url: string;
	    description: string;
	    is_enabled: number;
	    sort_order: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateMcpServerRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.command = source["command"];
	        this.args = source["args"];
	        this.env = source["env"];
	        this.url = source["url"];
	        this.description = source["description"];
	        this.is_enabled = source["is_enabled"];
	        this.sort_order = source["sort_order"];
	    }
	}
	export class CreateProfileRequest {
	    name: string;
	    provider_id: number;
	    model: string;
	    system_prompt: string;
	    temperature: number;
	    max_tokens: number;
	    extra_config: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateProfileRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.provider_id = source["provider_id"];
	        this.model = source["model"];
	        this.system_prompt = source["system_prompt"];
	        this.temperature = source["temperature"];
	        this.max_tokens = source["max_tokens"];
	        this.extra_config = source["extra_config"];
	    }
	}
	export class CreateProviderRequest {
	    name: string;
	    type: string;
	    api_url: string;
	    api_key: string;
	    models: string;
	    sort_order: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateProviderRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.api_url = source["api_url"];
	        this.api_key = source["api_key"];
	        this.models = source["models"];
	        this.sort_order = source["sort_order"];
	    }
	}
	export class CreateProxyRequest {
	    name: string;
	    type: string;
	    host: string;
	    port: number;
	    username: string;
	    password: string;
	    no_proxy: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateProxyRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.no_proxy = source["no_proxy"];
	    }
	}
	export class CreateSkillRequest {
	    name: string;
	    category: string;
	    trigger: string;
	    content: string;
	    variables: string;
	    sort_order: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateSkillRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.category = source["category"];
	        this.trigger = source["trigger"];
	        this.content = source["content"];
	        this.variables = source["variables"];
	        this.sort_order = source["sort_order"];
	    }
	}
	export class GetCliSessionsRequest {
	    cli_tool_key: string;
	    project: string;
	    limit: number;
	
	    static createFrom(source: any = {}) {
	        return new GetCliSessionsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cli_tool_key = source["cli_tool_key"];
	        this.project = source["project"];
	        this.limit = source["limit"];
	    }
	}
	export class GetSessionsRequest {
	    cli_tool_key: string;
	    page: number;
	    page_size: number;
	
	    static createFrom(source: any = {}) {
	        return new GetSessionsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cli_tool_key = source["cli_tool_key"];
	        this.page = source["page"];
	        this.page_size = source["page_size"];
	    }
	}
	export class LaunchRequest {
	    cli_tool_key: string;
	    profile_id: number;
	    proxy_id?: number;
	    mcp_server_ids: number[];
	    skill_ids: number[];
	    skill_vars: Record<string, string>;
	    extra_args: string[];
	    terminal: string;
	    working_dir: string;
	
	    static createFrom(source: any = {}) {
	        return new LaunchRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cli_tool_key = source["cli_tool_key"];
	        this.profile_id = source["profile_id"];
	        this.proxy_id = source["proxy_id"];
	        this.mcp_server_ids = source["mcp_server_ids"];
	        this.skill_ids = source["skill_ids"];
	        this.skill_vars = source["skill_vars"];
	        this.extra_args = source["extra_args"];
	        this.terminal = source["terminal"];
	        this.working_dir = source["working_dir"];
	    }
	}
	export class SetActiveConfigRequest {
	    cli_tool_key: string;
	    profile_id?: number;
	    proxy_id?: number;
	
	    static createFrom(source: any = {}) {
	        return new SetActiveConfigRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cli_tool_key = source["cli_tool_key"];
	        this.profile_id = source["profile_id"];
	        this.proxy_id = source["proxy_id"];
	    }
	}
	export class UpdateMcpServerRequest {
	    name: string;
	    type: string;
	    command: string;
	    args: string;
	    env: string;
	    url: string;
	    description: string;
	    is_enabled: number;
	    sort_order: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateMcpServerRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.command = source["command"];
	        this.args = source["args"];
	        this.env = source["env"];
	        this.url = source["url"];
	        this.description = source["description"];
	        this.is_enabled = source["is_enabled"];
	        this.sort_order = source["sort_order"];
	    }
	}
	export class UpdateProfileRequest {
	    name: string;
	    provider_id: number;
	    model: string;
	    system_prompt: string;
	    temperature: number;
	    max_tokens: number;
	    extra_config: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProfileRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.provider_id = source["provider_id"];
	        this.model = source["model"];
	        this.system_prompt = source["system_prompt"];
	        this.temperature = source["temperature"];
	        this.max_tokens = source["max_tokens"];
	        this.extra_config = source["extra_config"];
	    }
	}
	export class UpdateProviderRequest {
	    name: string;
	    type: string;
	    api_url: string;
	    api_key: string;
	    models: string;
	    sort_order: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProviderRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.api_url = source["api_url"];
	        this.api_key = source["api_key"];
	        this.models = source["models"];
	        this.sort_order = source["sort_order"];
	    }
	}
	export class UpdateProxyRequest {
	    name: string;
	    type: string;
	    host: string;
	    port: number;
	    username: string;
	    password: string;
	    no_proxy: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProxyRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.no_proxy = source["no_proxy"];
	    }
	}
	export class UpdateSkillRequest {
	    name: string;
	    category: string;
	    trigger: string;
	    content: string;
	    variables: string;
	    sort_order: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateSkillRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.category = source["category"];
	        this.trigger = source["trigger"];
	        this.content = source["content"];
	        this.variables = source["variables"];
	        this.sort_order = source["sort_order"];
	    }
	}

}

export namespace models {
	
	export class CliTool {
	    id: number;
	    name: string;
	    key: string;
	    executable: string;
	    config_path: string;
	    preferred_terminal: string;
	    is_installed: number;
	    is_enabled: number;
	    sort_order: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new CliTool(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.key = source["key"];
	        this.executable = source["executable"];
	        this.config_path = source["config_path"];
	        this.preferred_terminal = source["preferred_terminal"];
	        this.is_installed = source["is_installed"];
	        this.is_enabled = source["is_enabled"];
	        this.sort_order = source["sort_order"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class McpServer {
	    id: number;
	    name: string;
	    type: string;
	    command: string;
	    args: string;
	    env: string;
	    url: string;
	    description: string;
	    is_enabled: number;
	    sort_order: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new McpServer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.command = source["command"];
	        this.args = source["args"];
	        this.env = source["env"];
	        this.url = source["url"];
	        this.description = source["description"];
	        this.is_enabled = source["is_enabled"];
	        this.sort_order = source["sort_order"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Profile {
	    id: number;
	    name: string;
	    provider_id: number;
	    model: string;
	    system_prompt: string;
	    temperature: number;
	    max_tokens: number;
	    extra_config: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.provider_id = source["provider_id"];
	        this.model = source["model"];
	        this.system_prompt = source["system_prompt"];
	        this.temperature = source["temperature"];
	        this.max_tokens = source["max_tokens"];
	        this.extra_config = source["extra_config"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Provider {
	    id: number;
	    name: string;
	    type: string;
	    api_url: string;
	    api_key: string;
	    models: string;
	    sort_order: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Provider(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.api_url = source["api_url"];
	        this.api_key = source["api_key"];
	        this.models = source["models"];
	        this.sort_order = source["sort_order"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Proxy {
	    id: number;
	    name: string;
	    type: string;
	    host: string;
	    port: number;
	    username: string;
	    password: string;
	    no_proxy: string;
	    is_active: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Proxy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.no_proxy = source["no_proxy"];
	        this.is_active = source["is_active"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Session {
	    id: number;
	    cli_tool_id: number;
	    profile_id?: number;
	    proxy_id?: number;
	    terminal: string;
	    working_dir: string;
	    extra_args: string;
	    status: string;
	    pid: number;
	    // Go type: time
	    started_at: any;
	    // Go type: time
	    ended_at?: any;
	
	    static createFrom(source: any = {}) {
	        return new Session(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.cli_tool_id = source["cli_tool_id"];
	        this.profile_id = source["profile_id"];
	        this.proxy_id = source["proxy_id"];
	        this.terminal = source["terminal"];
	        this.working_dir = source["working_dir"];
	        this.extra_args = source["extra_args"];
	        this.status = source["status"];
	        this.pid = source["pid"];
	        this.started_at = this.convertValues(source["started_at"], null);
	        this.ended_at = this.convertValues(source["ended_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Skill {
	    id: number;
	    name: string;
	    category: string;
	    trigger: string;
	    content: string;
	    variables: string;
	    is_builtin: number;
	    sort_order: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Skill(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.category = source["category"];
	        this.trigger = source["trigger"];
	        this.content = source["content"];
	        this.variables = source["variables"];
	        this.is_builtin = source["is_builtin"];
	        this.sort_order = source["sort_order"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace service {
	
	export class ActiveConfig {
	    profile_id?: number;
	    proxy_id?: number;
	
	    static createFrom(source: any = {}) {
	        return new ActiveConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.profile_id = source["profile_id"];
	        this.proxy_id = source["proxy_id"];
	    }
	}
	export class CliSession {
	    session_id: string;
	    cli_tool_key: string;
	    project: string;
	    project_dir: string;
	    slug: string;
	    first_message: string;
	    message_count: number;
	    user_count: number;
	    assistant_count: number;
	    model: string;
	    started_at: string;
	    last_active_at: string;
	
	    static createFrom(source: any = {}) {
	        return new CliSession(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.session_id = source["session_id"];
	        this.cli_tool_key = source["cli_tool_key"];
	        this.project = source["project"];
	        this.project_dir = source["project_dir"];
	        this.slug = source["slug"];
	        this.first_message = source["first_message"];
	        this.message_count = source["message_count"];
	        this.user_count = source["user_count"];
	        this.assistant_count = source["assistant_count"];
	        this.model = source["model"];
	        this.started_at = source["started_at"];
	        this.last_active_at = source["last_active_at"];
	    }
	}
	export class CliSessionMessage {
	    type: string;
	    content: string;
	    timestamp: string;
	    model?: string;
	    tokens_in?: number;
	    tokens_out?: number;
	    uuid: string;
	
	    static createFrom(source: any = {}) {
	        return new CliSessionMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.content = source["content"];
	        this.timestamp = source["timestamp"];
	        this.model = source["model"];
	        this.tokens_in = source["tokens_in"];
	        this.tokens_out = source["tokens_out"];
	        this.uuid = source["uuid"];
	    }
	}
	export class CliSessionProject {
	    dir_name: string;
	    path: string;
	    session_count: number;
	
	    static createFrom(source: any = {}) {
	        return new CliSessionProject(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dir_name = source["dir_name"];
	        this.path = source["path"];
	        this.session_count = source["session_count"];
	    }
	}
	export class DetectResult {
	    key: string;
	    is_installed: boolean;
	    executable: string;
	
	    static createFrom(source: any = {}) {
	        return new DetectResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.is_installed = source["is_installed"];
	        this.executable = source["executable"];
	    }
	}
	export class TerminalInfo {
	    id: string;
	    name: string;
	    is_available: boolean;
	
	    static createFrom(source: any = {}) {
	        return new TerminalInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.is_available = source["is_available"];
	    }
	}
	export class TestResult {
	    success: boolean;
	    message: string;
	    latency_ms: number;
	
	    static createFrom(source: any = {}) {
	        return new TestResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.latency_ms = source["latency_ms"];
	    }
	}

}

