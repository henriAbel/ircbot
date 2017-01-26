import { Component, OnInit } from '@angular/core';
import { ApiService } from '../shared/'

@Component({
  selector: 'irc-stat',
  templateUrl: './stat.component.html',
})
export class StatComponent implements OnInit {

  private colors = ['#1BE7FF', '#6EEB83', '#E4FF1A', '#E8AA14', '#FF5714', '#50514F', '#F25F5C', '#247BA0', '#70C1B3'];
  private type: string;
  private options: any;
  private linkData: any;
  private userData: any;
  private duplicateData: any; 

  constructor(private apiService: ApiService) { 
    this.type = 'doughnut'
    this.options = {
      responsive: true,
      maintainAspectRatio: false
    };
    this.linkData = this.initDataSets();
    this.userData = this.initDataSets();
    this.duplicateData = this.initDataSets();
  }

  capitalizeFirstLetter(s: string) {
    return s.charAt(0).toUpperCase() + s.slice(1);
  }

  shuffle(o) {
    for (let j, x, i = o.length; i; j = Math.floor(Math.random() * i), x = o[--i], o[i] = o[j], o[j] = x);
  }
  getColor(i) {
    if (this.colors.length <= i) {
      let start = i - this.colors.length * Math.floor(i / this.colors.length);
      return this.colors[start];
    }
    return this.colors[i];
  }

  ngOnInit() {
    this.shuffle(this.colors);
    this.apiService.getStats().subscribe(data => {
      for (let i = 0; i < data.GroupLink.length; i++) {
        let e = data.GroupLink[i]
        this.linkData.labels.push(this.capitalizeFirstLetter(e.Type));
        this.linkData.datasets[0].data.push(e.Count);
        this.linkData.datasets[0].backgroundColor.push(this.getColor(i));
      }
      for (let i = 0; i < data.GroupUser.length; i++) {
        let e = data.GroupUser[i];
        this.userData.labels.push(this.capitalizeFirstLetter(e.User_name));
        this.userData.datasets[0].data.push(e.Count);
        this.userData.datasets[0].backgroundColor.push(this.getColor(i));
      }
      for (let i = 0; i < data.Duplicates.length; i++) {
        let e = data.Duplicates[i];
        this.duplicateData.labels.push(this.capitalizeFirstLetter(e.User_name));
        this.duplicateData.datasets[0].data.push(e.Count);
        this.duplicateData.datasets[0].backgroundColor.push(this.getColor(i));
      }
    })
  }

  initDataSets(): any {
    return {
      labels: [],
      datasets: [{
        data: [],
        backgroundColor: []
      }]
    }
  }

}
