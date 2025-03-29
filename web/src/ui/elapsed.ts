
/**
 * ElapsedTime class handles UI components to show relative time delay with auto-update
 * One class instance can manage several timers components
 */
export class ElapsedTime {
	timers: JQuery[];
    timer_id: number;

	constructor() {
		this.timers = [];
		this.timer_id = setInterval( () => {
			this.update();
		}, 10000);
	}
	
	/**
	* @param Date time
	 */
	create(time: Date) {
		const $e = $('<span class="date"/>');
		$e.attr('title', time.toLocaleString());
		const d = time.getTime() / 1000;
		$e.data('time', d);
		$e.text(human_delay(delay_seconds(d)));
		this.timers.push($e);
		return $e;
	}
	
	update() {
		const now = now_seconds();
		this.timers.forEach(($t) => {
			const time = $t.data('time');
			const d = delay_seconds(time, now);
			$t.text(human_delay(d));
		});
	}	
}

const now_seconds = ()=>{
	return Date.now() / 1000;
}

const delay_seconds = (time: number, now?:number) =>{
	if(!now) {
		now = now_seconds();
	}
	return now - time;
}

const days_seconds = 24 * 60 * 60;
const hours_seconds = 60 * 60;

/**
 * @param d delay in second
 * @returns 
 */
export const human_delay = (d: number) => {
	if(d > days_seconds) {
		return (~~(d/days_seconds)) + "d";
	}
	if(d > hours_seconds) {
		return (~~(d/(hours_seconds))) + "h";
	}
	if(d > 60) {
		return (~~(d / 60)) + "m";
	}
	return d + 's';
}
