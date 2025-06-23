# Comprehensive AML (Anti-Money Laundering) Rules and Checks

This document outlines the comprehensive AML rules and checks implemented in the accounting system, covering the most common and essential AML detection patterns used in financial institutions worldwide.

## Overview

The AML system implements multiple regulatory frameworks and detection methodologies to identify suspicious financial activities and ensure compliance with global anti-money laundering standards.

## Regulatory Frameworks Supported

### 1. BSA (Bank Secrecy Act) - United States
- **CTR (Currency Transaction Reports)**: Reports for cash transactions over $10,000
- **SAR (Suspicious Activity Reports)**: Reports for suspicious transactions over $5,000
- **Structuring Detection**: Identifies attempts to avoid reporting thresholds

### 2. FATF (Financial Action Task Force) - International
- **High-Risk Jurisdictions**: Monitoring transactions from/to high-risk countries
- **Politically Exposed Persons (PEP)**: Enhanced screening for politically exposed individuals
- **Trade-Based Money Laundering**: Detection of suspicious trade transactions

### 3. AMLD (Anti-Money Laundering Directive) - European Union
- **Enhanced Due Diligence**: Stricter verification for high-risk customers
- **Beneficial Ownership**: Identification of ultimate beneficial owners
- **Cross-Border Monitoring**: Enhanced monitoring of cross-border transactions

### 4. FinCEN (Financial Crimes Enforcement Network) - United States
- **Cross-Border Payments**: Monitoring international wire transfers
- **Virtual Currency**: Cryptocurrency transaction monitoring
- **Shell Company Detection**: Identification of shell company indicators

### 5. OFAC (Office of Foreign Assets Control) - United States
- **Sanctions Screening**: Real-time screening against sanctions lists
- **Blocked Persons**: Identification of sanctioned individuals and entities

## Core AML Rule Types

### Transaction-Based Rules

#### 1. Currency Transaction Report (CTR)
- **Threshold**: $10,000 USD
- **Purpose**: Report large cash transactions to regulatory authorities
- **Trigger**: Single transaction or daily aggregate cash amounts
- **Risk Level**: Medium to High

#### 2. Suspicious Activity Report (SAR)
- **Threshold**: $5,000 USD minimum (varies by jurisdiction)
- **Purpose**: Report potentially suspicious transactions
- **Factors**: Unusual patterns, lack of economic purpose, customer behavior
- **Risk Level**: High

#### 3. Structuring Detection
- **Pattern**: Multiple transactions designed to avoid reporting thresholds
- **Threshold**: Transactions just under $10,000 (typically 95-99% of threshold)
- **Time Window**: Short periods (hours to days)
- **Risk Level**: High

#### 4. Just Under Threshold
- **Pattern**: Transactions consistently just below reporting limits
- **Detection**: Amounts within 5% of thresholds, occurring 3+ times in 7 days
- **Risk Level**: High

#### 5. Round Amount Detection
- **Pattern**: Unusual frequency of round number transactions
- **Examples**: $1,000.00, $5,000.00, $10,000.00
- **Threshold**: 10+ round amounts in a period
- **Risk Level**: Medium

#### 6. Rapid Movement
- **Pattern**: Quick succession of transactions between accounts
- **Detection**: Multiple transactions within hours
- **Purpose**: Obscure money trail through layering
- **Risk Level**: Medium to High

#### 7. Layering
- **Pattern**: Complex transaction chains to obscure source
- **Detection**: 10+ intermediary accounts in money trail
- **Purpose**: Separate illicit funds from their source
- **Risk Level**: High

#### 8. Smurfing
- **Pattern**: Breaking large amounts into smaller transactions
- **Detection**: Multiple small transactions by related parties
- **Purpose**: Avoid detection and reporting thresholds
- **Risk Level**: High

### Customer-Based Rules

#### 9. Know Your Customer (KYC)
- **Requirement**: Customer identification and verification
- **Components**: Identity, address, business purpose verification
- **Frequency**: At onboarding and periodic updates
- **Risk Level**: Varies

#### 10. Customer Due Diligence (CDD)
- **Requirement**: Understanding customer's business and risk profile
- **Components**: Expected transaction patterns, source of funds
- **Updates**: Triggered by significant changes
- **Risk Level**: Medium

#### 11. Enhanced Due Diligence (EDD)
- **Trigger**: High-risk customers, PEPs, high-risk jurisdictions
- **Requirements**: Additional verification, ongoing monitoring
- **Frequency**: More frequent reviews and updates
- **Risk Level**: High

#### 12. Politically Exposed Persons (PEP)
- **Definition**: Individuals in prominent public positions
- **Screening**: Against PEP databases and lists
- **Monitoring**: Enhanced ongoing monitoring required
- **Risk Level**: High

### Geographic Rules

#### 13. High-Risk Jurisdictions
- **Countries**: FATF blacklisted and graylisted countries
- **Examples**: Afghanistan, Iran, North Korea, Myanmar, Syria
- **Monitoring**: Enhanced scrutiny for transactions to/from these countries
- **Risk Level**: High to Critical

#### 14. Unexpected Geography
- **Pattern**: Transactions from unexpected geographical locations
- **Detection**: Customer activity outside normal geographic patterns
- **Factors**: IP address, transaction location vs. customer address
- **Risk Level**: Medium to High

#### 15. NCCT (Non-Cooperative Countries and Territories)
- **Definition**: Countries with inadequate AML/CFT frameworks
- **Monitoring**: Enhanced monitoring of transactions
- **Updates**: Based on FATF recommendations
- **Risk Level**: High

### Pattern-Based Rules

#### 16. Velocity Monitoring
- **Pattern**: Unusual speed of transactions
- **Detection**: Transaction frequency exceeding normal patterns
- **Thresholds**: Customer-specific baselines
- **Risk Level**: Medium

#### 17. High Frequency
- **Pattern**: Unusually high number of transactions
- **Threshold**: 50+ transactions per day
- **Detection**: Account-level monitoring
- **Risk Level**: High

#### 18. Cash Intensive Activity
- **Pattern**: Unusually high percentage of cash transactions
- **Threshold**: 80% cash transactions with $50,000+ volume
- **Time Window**: 30 days
- **Risk Level**: High

#### 19. Unusual Timing
- **Pattern**: Transactions at unusual times
- **Detection**: Night time (10 PM - 6 AM), weekends, holidays
- **Minimum**: $1,000+ transactions
- **Risk Level**: Medium

#### 20. Account Dormancy Reactivation
- **Pattern**: Sudden activity in dormant accounts
- **Dormancy Period**: 90+ days of inactivity
- **Reactivation Threshold**: $5,000+ transaction
- **Risk Level**: Medium

### Advanced Detection Rules

#### 21. Wire Stripping
- **Pattern**: Removal of wire transfer information
- **Detection**: Missing or altered wire details
- **Minimum**: $3,000+ transactions
- **Risk Level**: High

#### 22. Third-Party Check Deposits
- **Pattern**: Deposits of checks from third parties
- **Threshold**: 5+ per month, $500+ minimum
- **Risk**: Check fraud, money laundering
- **Risk Level**: Medium

#### 23. Cryptocurrency Transactions
- **Pattern**: Virtual currency related activities
- **Monitoring**: All cryptocurrency transactions $1,000+
- **Risks**: Anonymity, cross-border movement
- **Risk Level**: Medium

#### 24. Shell Company Indicators
- **Pattern**: Characteristics of shell companies
- **Indicators**: Minimal operations, complex ownership
- **Detection**: Transaction complexity analysis
- **Risk Level**: High

#### 25. Trade-Based Money Laundering
- **Pattern**: Manipulation of trade transactions
- **Detection**: Price variance >20% from market rates
- **Volume**: $10,000+ transactions
- **Risk Level**: High

#### 26. Circular Transfers
- **Pattern**: Money returning to originator
- **Detection**: Circular flow analysis
- **Purpose**: Create appearance of legitimate activity
- **Risk Level**: High

#### 27. Concentration Risk
- **Pattern**: High concentration of activity in specific accounts
- **Detection**: Account-level concentration analysis
- **Thresholds**: Percentage-based limits
- **Risk Level**: Medium

## Risk Scoring

### Risk Levels
- **LOW**: Minimal risk, standard monitoring
- **MEDIUM**: Moderate risk, enhanced monitoring
- **HIGH**: Significant risk, investigation required
- **CRITICAL**: Immediate attention, potential SAR filing

### Risk Factors
1. **Transaction Amount**: Higher amounts = higher risk
2. **Customer Risk Profile**: PEP, high-risk jurisdiction
3. **Transaction Pattern**: Frequency, timing, routing
4. **Geographic Factors**: High-risk countries, unexpected locations
5. **Product Risk**: Cash, wire transfers, cryptocurrency

## Compliance Metrics

### Key Performance Indicators (KPIs)
- **False Positive Rate**: Target <10%
- **Alert Resolution Time**: Target <48 hours
- **CTR Filing Rate**: Percentage of cash transactions reported
- **SAR Filing Rate**: Percentage of suspicious transactions reported
- **Compliance Score**: Overall AML program effectiveness (0-100)

### Monitoring Dashboard
- Real-time alert monitoring
- Trend analysis (30-day rolling)
- Customer risk summaries
- Regulatory reporting status
- Investigation workflow management

## Implementation Guidelines

### Setup Process
1. **Framework Selection**: Choose applicable regulatory frameworks
2. **Rule Configuration**: Set thresholds based on business risk appetite
3. **Customer Onboarding**: Implement KYC/CDD procedures
4. **Ongoing Monitoring**: Establish transaction monitoring processes
5. **Alert Management**: Create investigation and escalation procedures

### Best Practices
1. **Regular Rule Tuning**: Adjust thresholds to minimize false positives
2. **Staff Training**: Regular AML training for all relevant personnel
3. **Documentation**: Maintain detailed records of decisions and investigations
4. **Technology Updates**: Keep screening lists and detection rules current
5. **Independent Testing**: Regular independent validation of AML program

### Integration Points
- **Transaction Processing**: Real-time monitoring during transaction processing
- **Customer Onboarding**: KYC/CDD checks during account opening
- **Periodic Reviews**: Scheduled customer risk assessments
- **Regulatory Reporting**: Automated CTR/SAR generation and filing
- **Case Management**: Investigation workflow and documentation

## Regulatory Reporting

### Required Reports
- **CTRs**: Currency Transaction Reports for cash transactions >$10,000
- **SARs**: Suspicious Activity Reports for suspicious patterns
- **CMIR**: Currency and Monetary Instrument Reports for border crossings
- **FBAR**: Foreign Bank Account Reports for offshore accounts

### Filing Timelines
- **CTRs**: 15 days after transaction
- **SARs**: 30 days after detection (60 days if no suspect identified)
- **CMIR**: At time of transport
- **FBAR**: Annual filing by April 15

This comprehensive AML framework ensures robust detection of money laundering activities while maintaining compliance with global regulatory requirements.
