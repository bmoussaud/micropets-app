import { BrowserModule } from '@angular/platform-browser';
import { NgModule, APP_INITIALIZER } from '@angular/core';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { PetsComponent } from './pets/pets.component';
import { HttpClientModule} from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MaterialModule } from './material.module';
import { ConfigAssetLoaderService} from './config-asset-loader.service';

@NgModule({
  declarations: [
    AppComponent,
    PetsComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MaterialModule
  ],
  providers: [{
    provide: APP_INITIALIZER,
    useFactory: (configService: ConfigAssetLoaderService) => () => configService.loadConfigurations().toPromise(),
    deps: [ConfigAssetLoaderService],
    multi: true
  }],
  bootstrap: [AppComponent]
})
export class AppModule { }
