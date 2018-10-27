import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { UsersComponent } from './users/users.component';
import { UserComponent } from './user/user.component';
import { AppComponent } from './app.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { EditProfileComponent } from './edit-profile/edit-profile.component';

export const routes: Routes = [
    { path: '', component: DashboardComponent, pathMatch: 'full' },
    { path: 'users', component: UsersComponent },
    { path: 'user/:id', component: UserComponent },
    { path: 'edit-profile', component: EditProfileComponent },
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutingModule { }
