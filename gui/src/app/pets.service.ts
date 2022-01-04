import { Injectable } from '@angular/core';
import { forkJoin, Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';


export interface PetsData {
  Total: number;
  Hostname: string;
  Hostnames: (HostnamesEntity)[]
  Pets: (PetsEntity)[]
}
export interface HostnamesEntity {
  Service: string;
  Hostname: string;
}
export interface PetsEntity {
  Index: number;
  Name: string;
  Type: string;
  Kind: string;
  Age: number;
  URL: string;
  Hostname: string;
  URI: string;
}


@Injectable({
  providedIn: 'root'
})
export class PetsService {

  constructor(private http: HttpClient) { }

  public getPetsData(url: string): Observable<PetsData> {
    if (url == "/") {
      return this.http.get<PetsData>("/pets");
    } else {
      return this.http.get<PetsData>(url + "/pets");
    }
  }

  public getPet(url: string): Observable<PetsEntity> {
    return this.http.get<PetsEntity>(url);
  }

  public getPets(urls: string[]): Observable<PetsEntity[]> {
    // firstly, start out with an array of observable arrays
    const observables: Observable<PetsEntity>[] = urls.map(url => this.getPet(url));
    // run all observables in parallel with forkJoin
    return forkJoin(observables).pipe(
      // now map the array of arrays to a flattened array
      map(pets => pets)
    );
  }
}
