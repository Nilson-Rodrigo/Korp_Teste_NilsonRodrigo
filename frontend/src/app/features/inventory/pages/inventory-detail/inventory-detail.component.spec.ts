import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { provideRouter } from '@angular/router';
import { InventoryDetailComponent } from './inventory-detail.component';

describe('InventoryDetailComponent', () => {
  let component: InventoryDetailComponent;
  let fixture: ComponentFixture<InventoryDetailComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [InventoryDetailComponent],
      providers: [
        provideHttpClient(),
        provideHttpClientTesting(),
        provideRouter([]),
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(InventoryDetailComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
