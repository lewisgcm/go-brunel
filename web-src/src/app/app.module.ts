import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientModule } from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

import { UserModule } from './module/user';
import { RepositoryModule } from './module/repository';
import { JobModule } from './module/job';
import { SharedModule } from './module/shared';

@NgModule({
    declarations: [
        AppComponent
    ],
    imports: [
        BrowserModule,
        BrowserAnimationsModule,
        HttpClientModule,
        AppRoutingModule,
        SharedModule,
        RepositoryModule,
        UserModule,
        JobModule,
    ],
    providers: [],
    bootstrap: [AppComponent]
})
export class AppModule { }
