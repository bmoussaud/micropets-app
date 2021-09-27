import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Location } from '@angular/common';
import { map } from 'rxjs/operators'
import { MatTableDataSource } from '@angular/material/table';
import { ConfigAssetLoaderService, Configuration } from '../config-asset-loader.service';


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
  Name: string;
  Type: string;
  Kind: string;
  Age: number;
  URL: string;
  Hostname: string;
}


@Component({
  selector: 'app-pets',
  templateUrl: './pets.component.html',
  styleUrls: ['./pets.component.css']
})
export class PetsComponent implements OnInit {

  public pets: PetsEntity[] = []
  public hostnames: any[] = []
  public env: string = ""
  public hostnamesstr: string = ""
  public env_color: string = "pink"
  public dataSource!: MatTableDataSource<any>;
  public dataSourceHostnames!: MatTableDataSource<any>;

  public config!: Configuration

  displayedColumns = ['name', 'kind', 'age', 'pic']

  constructor(private http: HttpClient, private router: Router, private location: Location, private configService: ConfigAssetLoaderService) {    
    
    this.configService.loadConfigurations().subscribe((data: Configuration) => this.config = {
      petServiceUrl: data.petServiceUrl,
      stage: data.stage,
      stage_color: data.stage_color,
    });
  }

  ngOnInit() {
    this.location.subscribe(() => {
      this.refresh();
    });
    this.refresh();
  }

  private refresh() {
    //console.log("------------------- refresh")
    //console.log(this.config.petServiceUrl)
    this.http.get<PetsData>(this.config.petServiceUrl)
      .pipe(map(result => result))
      .subscribe(result => {
        this.pets = result['Pets'];
        this.hostnames = result['Hostnames'];
        var h: string[] = new Array(4)
        for (let index = 0; index < this.hostnames.length; index++) {
          const element = this.hostnames[index];
          h[index] = element.Hostname
        }
        this.hostnamesstr = h.join(", ")

        this.env = this.config.stage;
        this.env_color = this.config.stage_color;
        this.dataSource = new MatTableDataSource(this.pets);
        this.dataSourceHostnames = new MatTableDataSource(this.hostnames)
      });
  }


}
