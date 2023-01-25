import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { shareReplay } from 'rxjs/operators';

export interface Configuration {
  petServiceUrl: string;
  stage: string;
  stage_color: string;
  load_one_by_one: string
}

@Injectable({ providedIn: 'root' })
export class ConfigAssetLoaderService {

  private readonly CONFIG_URL = 'assets/app-config/config.json';
  private configuration$!: Observable<Configuration>;

  constructor(private http: HttpClient) {
  }

  public loadConfigurations(): any {
    if (!this.configuration$) {
      this.configuration$ = this.http.get<Configuration>(this.CONFIG_URL).pipe(
        shareReplay(1)
      );
    }
    return this.configuration$;
  }

}
