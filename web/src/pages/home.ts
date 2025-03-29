import { Page } from "./page";

export class Home extends Page {

    constructor() {
        super('home');
    }

    render(): void {
        const $list = $('#home-namepaces');
        $list.empty();
    }
}