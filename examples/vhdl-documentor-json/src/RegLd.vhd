library IEEE;

use IEEE.std_logic_1164.all;

entity RegLd is

  port(
       iCLK             : in std_logic;
       iD               : in integer;
       iLd              : in integer;
       oQ               : out integer
     );

end RegLd;

architecture behavior of RegLd is
  signal s_Q : integer;
begin
  

  process(iCLK, iLd, iD)
  begin
    if rising_edge(iCLK) then
      if (iLd = 1) then
        s_Q <= iD;
      else
        s_Q <= s_Q;
      end if;
    end if;
  end process;

  oQ <= s_Q; -- connect internal storage signal with final output
  
end behavior;
