export namespace database {
	
	export class Camera {
	    id: number;
	    name: string;
	    rtsp_url: string;
	    is_active: boolean;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new Camera(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.rtsp_url = source["rtsp_url"];
	        this.is_active = source["is_active"];
	        this.status = source["status"];
	    }
	}
	export class Event {
	    id: number;
	    zone_id: number;
	    event_type: string;
	    person_count: number;
	    duration_seconds: number;
	    timestamp: string;
	    synced: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Event(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.zone_id = source["zone_id"];
	        this.event_type = source["event_type"];
	        this.person_count = source["person_count"];
	        this.duration_seconds = source["duration_seconds"];
	        this.timestamp = source["timestamp"];
	        this.synced = source["synced"];
	    }
	}
	export class Zone {
	    id: number;
	    camera_id: number;
	    name: string;
	    x1: number;
	    y1: number;
	    x2: number;
	    y2: number;
	    threshold_seconds: number;
	
	    static createFrom(source: any = {}) {
	        return new Zone(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.camera_id = source["camera_id"];
	        this.name = source["name"];
	        this.x1 = source["x1"];
	        this.y1 = source["y1"];
	        this.x2 = source["x2"];
	        this.y2 = source["y2"];
	        this.threshold_seconds = source["threshold_seconds"];
	    }
	}

}

export namespace models {
	
	export class DiscoveredCamera {
	    name: string;
	    endpoint: string;
	
	    static createFrom(source: any = {}) {
	        return new DiscoveredCamera(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.endpoint = source["endpoint"];
	    }
	}

}

