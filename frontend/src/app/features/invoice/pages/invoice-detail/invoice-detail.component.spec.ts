import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { provideRouter } from '@angular/router';
import { InvoiceDetailComponent } from './invoice-detail.component';

describe('InvoiceDetailComponent', () => {
  let component: InvoiceDetailComponent;
  let fixture: ComponentFixture<InvoiceDetailComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [InvoiceDetailComponent],
      providers: [
        provideHttpClient(),
        provideHttpClientTesting(),
        provideRouter([]),
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(InvoiceDetailComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
