import { Component, Output, EventEmitter, Input, OnDestroy } from '@angular/core';
import { Repository } from '../../../shared';

@Component({
    selector: 'app-repository-list',
    templateUrl: './list.component.html',
    styleUrls: ['./list.component.scss'],
})
export class ListComponent implements OnDestroy {
    @Output() search: EventEmitter<string>;
    @Output() repositorySelect: EventEmitter<Repository>;
    @Input() repositories: Repository[];
    @Input() loading = false;

    constructor() {
        this.search = new EventEmitter();
        this.repositorySelect = new EventEmitter();
    }

    ngOnDestroy() {
        this.repositorySelect.complete();
        this.search.complete();
    }
}
