import { Component, Input, Output, EventEmitter } from '@angular/core';
import { SortDirection } from '@angular/material/sort';

import { Repository } from '../../../shared';

export interface TableEvent {
    filter: string;
    pageIndex: number;
    pageSize: number;
    sortColumn: string;
    sortOrder: SortDirection;
}

@Component({
    selector: 'app-repository-detail',
    templateUrl: './detail.component.html',
    styleUrls: ['./detail.component.scss'],
})
export class DetailComponent {
    @Input() repository: Repository;
    @Output() jobSelect: EventEmitter<string> = new EventEmitter();

    public onJobClick(e: string) {
        this.jobSelect.emit(e);
    }
}
